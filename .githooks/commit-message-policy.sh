#!/bin/sh

validate_commit_message() {
	message="$1"
	ticket='GYXT-35505'
	ticket_url='https://devops.aliyun.com/projex/req/GYXT-35505#'

	if ! printf '%s\n' "$message" | grep -Eq "(^|[^[:alnum:]_-])${ticket}([^[:alnum:]_-]|$)"; then
		printf '提交注释必须包含完整工单编号 %s。\n' "$ticket" >&2
		return 1
	fi

	if ! printf '%s\n' "$message" | grep -Fq "$ticket_url"; then
		printf '提交注释必须包含云效工单链接：%s\n' "$ticket_url" >&2
		return 1
	fi
}
