#!/usr/bin/env python3
"""Local-only, bounded XLSX inventory conversion (Python standard library)."""
import argparse
import json
import pathlib
import posixpath
import re
import sys
import zipfile
import xml.etree.ElementTree as ET
from xml.sax.saxutils import escape

HEADERS = ['面料编码', '面料名称', '色号', '公司专属成品库存']
FIELDS = ['material_code', 'material_name', 'color_number', 'reserve_for_nsy']
NS = '{http://schemas.openxmlformats.org/spreadsheetml/2006/main}'
LIMIT = 16 * 1024 * 1024

def read(path):
    path = pathlib.Path(path)
    if path.suffix.lower() != '.xlsx':
        raise ValueError('仅支持 .xlsx；不支持 .xls 或其他格式')
    with zipfile.ZipFile(path) as z:
        infos = z.infolist()
        if len(infos) > 200 or sum(i.file_size for i in infos) > LIMIT:
            raise ValueError('XLSX expansion exceeds limit')
        def xml(name):
            data = z.read(name)
            if b'<!DOCTYPE' in data.upper() or b'<!ENTITY' in data.upper():
                raise ValueError('XML entities forbidden')
            return ET.fromstring(data)
        strings = []
        date_styles = set()
        if 'xl/styles.xml' in z.namelist():
            styles = xml('xl/styles.xml')
            custom = {int(e.attrib['numFmtId']): e.attrib.get('formatCode', '') for e in styles.findall(NS+'numFmts/'+NS+'numFmt')}
            for index, xf in enumerate(styles.findall(NS+'cellXfs/'+NS+'xf')):
                number_format = int(xf.attrib.get('numFmtId', '0'))
                fmt = custom.get(number_format, '')
                # Remove quoted/escaped literals and non-time bracket clauses.
                fmt = re.sub(r'"[^"]*"|\\.|_.|\*.', '', fmt)
                fmt = re.sub(r'\[(?![hms]+\])[^\]]*\]', '', fmt, flags=re.I)
                if number_format in set(range(14, 23)) | set(range(27, 37)) | set(range(45, 48)) | set(range(50, 59)) or re.search(r'[ymdhs]|AM/PM|A/P', fmt, re.I):
                    date_styles.add(index)
        if 'xl/sharedStrings.xml' in z.namelist():
            strings = [''.join(e.itertext()) for e in xml('xl/sharedStrings.xml').findall(NS+'si')]
        workbook = xml('xl/workbook.xml')
        sheets = workbook.findall(NS+'sheets/'+NS+'sheet')
        if not sheets:
            raise ValueError('Workbook has no worksheet')
        relation = sheets[0].attrib.get('{http://schemas.openxmlformats.org/officeDocument/2006/relationships}id')
        relationships = xml('xl/_rels/workbook.xml.rels')
        targets = [r for r in relationships if r.attrib.get('Id') == relation]
        if len(targets) != 1 or targets[0].attrib.get('TargetMode') == 'External' or not targets[0].attrib.get('Type', '').endswith('/worksheet'):
            raise ValueError('Invalid worksheet relationship')
        target = targets[0].attrib.get('Target', '')
        sheet_path = posixpath.normpath(target.lstrip('/') if target.startswith('/') else 'xl/'+target)
        if not sheet_path.startswith('xl/') or '..' in sheet_path.split('/'):
            raise ValueError('Invalid worksheet path')
        root = xml(sheet_path)
        rows = []
        for row in root.findall('.//'+NS+'row'):
            n = int(row.attrib['r'])
            if n > 10000 or n < 1:
                raise ValueError('Worksheet row limit')
            while len(rows) < n:
                rows.append(['']*5)
            for c in row.findall(NS+'c'):
                ref = c.attrib.get('r', '')
                m = re.fullmatch(r'([A-Z]+)[0-9]+', ref)
                if not m or m[1] not in 'ABCDE' or len(m[1]) != 1:
                    raise ValueError('Unexpected worksheet column')
                idx = ord(m[1])-65
                if c.find(NS+'f') is not None:
                    rows[n-1][idx] = {'formula': c.find(NS+'f').text or ''}
                    continue
                typ = c.attrib.get('t')
                value = c.findtext(NS+'v', '')
                if typ == 's':
                    value = strings[int(value)]
                elif typ == 'inlineStr':
                    value = ''.join(c.find(NS+'is').itertext())
                elif idx in (0, 2) and value != '':
                    value = {'nontext': value, 'type': typ or 'n'}
                elif idx == 3 and (typ in ('b', 'e', 'd') or (typ in (None, 'n') and int(c.attrib.get('s', '0')) in date_styles)):
                    value = {'nontext': value, 'type': typ or 'n', 'date_style': typ in (None, 'n')}
                rows[n-1][idx] = value
        return rows

def write(path, rows):
    path = pathlib.Path(path)
    if path.suffix.lower() != '.xlsx':
        raise ValueError('Output must be .xlsx')
    cells = []
    for n, row in enumerate(rows, 1):
        cs = []
        for i, value in enumerate(row):
            if isinstance(value, dict) and 'nontext' in value:
                typ = value.get('type', 'n')
                if typ not in ('n', 'b', 'e', 'd', 'str'):
                    raise ValueError('Invalid raw cell type')
                cs.append('<c r="%s%d" t="%s"%s><v>%s</v></c>' % (chr(65+i), n, typ, ' s="1"' if value.get('date_style') else '', escape(str(value['nontext']))))
                continue
            if isinstance(value, dict):
                value = '='+value['formula'] if 'formula' in value else value.get('nontext', '')
            cs.append('<c r="%s%d" t="inlineStr"><is><t xml:space="preserve">%s</t></is></c>' % (chr(65+i), n, escape(str(value))))
        cells.append('<row r="%d">%s</row>' % (n, ''.join(cs)))
    with path.open('xb') as f, zipfile.ZipFile(f, 'w', zipfile.ZIP_DEFLATED) as z:
        z.writestr('xl/styles.xml', '<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><cellXfs count="2"><xf numFmtId="0"/><xf numFmtId="14" applyNumberFormat="1"/></cellXfs></styleSheet>')
        z.writestr('[Content_Types].xml', '<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/xl/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"/><Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/><Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/></Types>')
        z.writestr('_rels/.rels', '<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/></Relationships>')
        z.writestr('xl/workbook.xml', '<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets><sheet name="库存" sheetId="1" r:id="rId1"/></sheets></workbook>')
        z.writestr('xl/_rels/workbook.xml.rels', '<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/><Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/></Relationships>')
        z.writestr('xl/worksheets/sheet1.xml', '<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>'+''.join(cells)+'</sheetData></worksheet>')

def parse(path):
    rows = read(path)
    if not rows or rows[0][:4] != HEADERS or rows[0][4:] not in ([], [''], ['错误原因']):
        raise ValueError('模板须为固定四列或追加错误原因的五列')
    result = {'items': [], 'row_mapping': [], 'raw_rows': [r[:4] for r in rows[1:]], 'errors': []}
    for n, row in enumerate(rows[1:], 2):
        row = row[:4]
        if all(v == '' for v in row):
            continue
        reasons = []
        if any(isinstance(v, dict) and 'formula' in v for v in row):
            reasons.append('不接受公式输入；请粘贴为值')
        if any(isinstance(row[i], dict) and 'nontext' in row[i] for i in (0, 2)):
            reasons.append('面料编码和色号必须为文本单元格；请按实际编码重新输入为文本，避免前导零丢失')
        if not row[0] or not row[2]:
            reasons.append('面料编码和色号必填')
        stock = str(row[3])
        if isinstance(row[3], dict) and 'nontext' in row[3]:
            reasons.append('库存不接受布尔、错误或日期单元格；请输入整数')
        if not re.fullmatch(r'[0-9]+', stock) or len(stock) > 10 or int(stock) > 2147483647:
            reasons.append('目标库存须为非负 int32 整数，空值不是0')
        if reasons:
            result['errors'].append({'excel_row': n, 'message': '；'.join(reasons)})
            continue
        item = dict(zip(FIELDS, row))
        item['reserve_for_nsy'] = int(stock)
        if not str(item['material_name']).strip():
            item['material_name'] = ''
        result['items'].append(item)
        result['row_mapping'].append(n)
    if not 1 <= len(result['row_mapping']) + len(result['errors']) <= 500:
        raise ValueError('每批必须1–500行；不得自动拆分')
    return result

def errors(envelope, result, output):
    messages = {}
    for error in result.get('errors', []):
        n = error.get('excel_row')
        if n is None:
            index = error.get('row_index')
            if type(index) is not int or not 1 <= index <= len(envelope['row_mapping']):
                raise ValueError('Invalid MES row_index')
            n = envelope['row_mapping'][index-1]
        messages.setdefault(n, []).append(str(error['message']))
    truncated = result.get('errors_truncated', False)
    if not messages and not truncated:
        raise ValueError('没有行错误；认证/权限/网络错误不生成错误文件')
    note = '错误详情未完整返回；须整批重新预检' if truncated else ''
    write(output, [HEADERS+['错误原因']] + [row+['；'.join(messages.get(n, [])+([note] if note else []))] for n, row in enumerate(envelope['raw_rows'], 2)])
    if truncated:
        print(note, file=sys.stderr)

def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('command', choices=['template', 'parse', 'export', 'errors'])
    p.add_argument('--input'); p.add_argument('--output', required=True)
    p.add_argument('--result'); p.add_argument('--errors-output')
    a = p.parse_args()
    if a.command == 'template':
        write(a.output, [HEADERS]); return
    if a.command == 'parse':
        data = parse(a.input)
        if data['errors']:
            if not a.errors_output:
                raise ValueError('--errors-output required on validation failure')
            errors(data, data, a.errors_output)
        with open(a.output, 'x', encoding='utf-8') as f:
            json.dump(data, f, ensure_ascii=False, indent=2)
        if data['errors']:
            raise ValueError('整批未提交；修正错误文件后重新解析')
        return
    with open(a.input, encoding='utf-8') as f:
        data = json.load(f)
    if a.command == 'export':
        write(a.output, [HEADERS]+[[item.get(k, '') for k in FIELDS] for item in data['items']])
    else:
        with open(a.result, encoding='utf-8') as f:
            errors(data, json.load(f), a.output)

if __name__ == '__main__':
    try:
        main()
    except (ValueError, OSError, KeyError, IndexError, ET.ParseError, zipfile.BadZipFile) as e:
        print(str(e), file=sys.stderr)
        sys.exit(1)
