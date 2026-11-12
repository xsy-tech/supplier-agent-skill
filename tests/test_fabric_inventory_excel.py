import importlib.util
import pathlib
import tempfile
import unittest
import zipfile

SCRIPT = pathlib.Path(__file__).resolve().parents[1] / 'scripts/fabric_inventory_excel.py'
spec = importlib.util.spec_from_file_location('excel', SCRIPT)
excel = importlib.util.module_from_spec(spec)
spec.loader.exec_module(excel)

class InventoryTest(unittest.TestCase):
    def test_roundtrip_and_mapping(self):
        with tempfile.TemporaryDirectory() as d:
            p = pathlib.Path(d)/'input.xlsx'
            excel.write(p, [excel.HEADERS, ['001', '', '01', 0], ['', '', '', ''], ['002', '=literal', '02', 12]])
            result = excel.parse(p)
            self.assertEqual(result['row_mapping'], [2, 4])
            self.assertEqual(result['items'][0]['material_code'], '001')
            self.assertEqual(result['items'][0]['reserve_for_nsy'], 0)
            out = pathlib.Path(d)/'errors.xlsx'
            excel.errors(result, {'errors':[{'row_index':2,'message':'bad'}]}, out)
            rows = excel.read(out)
            self.assertEqual(rows[3][4], 'bad')
            self.assertEqual(excel.parse(out)['items'], result['items'])
            with self.assertRaises(FileExistsError):
                excel.write(out, [])

    def test_invalid_stock(self):
        with tempfile.TemporaryDirectory() as d:
            p = pathlib.Path(d)/'x.xlsx'
            excel.write(p, [excel.HEADERS, ['001','','01',''], ['002','','02',2147483648]])
            self.assertEqual(len(excel.parse(p)['errors']), 2)

    def test_formula_and_nonrow_errors(self):
        with tempfile.TemporaryDirectory() as d:
            p = pathlib.Path(d)/'x.xlsx'
            excel.write(p, [excel.HEADERS, ['001','','01',1]])
            with zipfile.ZipFile(p) as z:
                parts = {n:z.read(n) for n in z.namelist()}
            parts['xl/worksheets/sheet1.xml'] = parts['xl/worksheets/sheet1.xml'].replace(b'<c r="D2" t="inlineStr"><is><t xml:space="preserve">1</t></is></c>', b'<c r="D2"><f>1+1</f><v>2</v></c>')
            with zipfile.ZipFile(p, 'w') as z:
                for name, value in parts.items():
                    z.writestr(name, value)
            data = excel.parse(p)
            self.assertIn('公式', data['errors'][0]['message'])
            with self.assertRaises(ValueError):
                excel.errors(data, {'code':'UNAUTHORIZED'}, pathlib.Path(d)/'auth.xlsx')
            out = pathlib.Path(d)/'truncated.xlsx'
            excel.errors(data, {'errors_truncated':True}, out)
            self.assertIn('错误详情未完整返回', excel.read(out)[1][4])

    def test_limits_and_formats(self):
        with tempfile.TemporaryDirectory() as d:
            with self.assertRaises(ValueError):
                excel.read(pathlib.Path(d)/'input.xls')
            p = pathlib.Path(d)/'large.xlsx'
            excel.write(p, [excel.HEADERS]+[['001','','01',0]]*501)
            with self.assertRaises(ValueError):
                excel.parse(p)

    def test_numeric_identifier_and_relationship_selection(self):
        with tempfile.TemporaryDirectory() as d:
            p = pathlib.Path(d)/'renamed.xlsx'
            excel.write(p, [excel.HEADERS, ['001','','01',0]])
            with zipfile.ZipFile(p) as z:
                parts = {n:z.read(n) for n in z.namelist()}
            sheet = parts.pop('xl/worksheets/sheet1.xml')
            parts['xl/worksheets/renamed.xml'] = sheet
            parts['xl/_rels/workbook.xml.rels'] = parts['xl/_rels/workbook.xml.rels'].replace(b'sheet1.xml', b'renamed.xml')
            def save():
                with zipfile.ZipFile(p, 'w') as z:
                    for name, value in parts.items():
                        z.writestr(name, value)
            save()
            self.assertEqual(excel.parse(p)['items'][0]['material_code'], '001')
            parts['xl/worksheets/renamed.xml'] = sheet.replace(b'<c r="A2" t="inlineStr"><is><t xml:space="preserve">001</t></is></c>', b'<c r="A2" s="1"><v>1</v></c>')
            parts['xl/styles.xml'] = b'<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><numFmts count="1"><numFmt numFmtId="164" formatCode="000"/></numFmts><cellXfs count="2"><xf numFmtId="0"/><xf numFmtId="164"/></cellXfs></styleSheet>'
            save()
            result = excel.parse(p)
            self.assertEqual(result['items'], [])
            self.assertIn('文本单元格', result['errors'][0]['message'])
            error_file = pathlib.Path(d)/'errors.xlsx'
            excel.errors(result, result, error_file)
            reparsed = excel.parse(error_file)
            self.assertEqual(reparsed['items'], [])
            self.assertIn('文本单元格', reparsed['errors'][0]['message'])
            parts['xl/worksheets/renamed.xml'] = sheet.replace(b'<c r="D2" t="inlineStr"><is><t xml:space="preserve">0</t></is></c>', b'<c r="D2" t="b"><v>1</v></c>')
            save()
            self.assertEqual(excel.parse(p)['items'], [])
            self.assertIn('布尔', excel.parse(p)['errors'][0]['message'])
            for numfmt, custom in [(14, ''), (164, '<numFmts count="1"><numFmt numFmtId="164" formatCode="yyyy-mm-dd"/></numFmts>')]:
                parts['xl/styles.xml'] = ('<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">'+custom+'<cellXfs count="2"><xf numFmtId="0"/><xf numFmtId="%d"/></cellXfs></styleSheet>' % numfmt).encode()
                parts['xl/worksheets/renamed.xml'] = sheet.replace(b'<c r="D2" t="inlineStr"><is><t xml:space="preserve">0</t></is></c>', b'<c r="D2" s="1"><v>45000</v></c>')
                save()
                self.assertEqual(excel.parse(p)['items'], [])
                self.assertIn('日期', excel.parse(p)['errors'][0]['message'])

if __name__ == '__main__':
    unittest.main()
