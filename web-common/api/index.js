import request from '@/web-common/axios';
export function getUser() {
  return request('GET', 'user');
}
