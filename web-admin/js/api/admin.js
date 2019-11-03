import request from '@/web-common/axios';

export function getSubmitStatistic() {
  return request('GET', 'admin/submitstatistic');
}
