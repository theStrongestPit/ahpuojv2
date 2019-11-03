import request from '@/web-common/axios';

export function toggleIssueStatus(id) {
  return request('PUT', `admin/issue/${id}/status`);
}

export function toggleReplyStatus(id) {
  return request('PUT', `admin/reply/${id}/status`);
}
