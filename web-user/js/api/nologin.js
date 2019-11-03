import request from '@/web-common/axios';

export function getLanguageList() {
  return request('GET', 'languages');
}
export function getRankList(page, perpage) {
  return request('GET', 'ranklist', {page: page, perpage: perpage});
}
export function getNewList(page, perpage) {
  return request('GET', 'news', {page: page, perpage: perpage});
}
export function getIssueList(id, page, perpage) {
  return request('GET', `problem/${id}/issues`, {page: page, perpage: perpage});
}
export function getIssue(id, page, perpage) {
  return request('GET', `/issue/${id}`, {page: page, perpage: perpage});
}
export function getProblem(id) {
  return request('GET', `problem/${id}`);
}
export function getProblemList(page, perpage, queryParam, level, tagId) {
  return request('GET', 'problems', {
    page: page,
    perpage: perpage,
    param: queryParam,
    level: level,
    tag_id: tagId
  });
}
export function getSolution(id) {
  return request('GET', `solution/${id}`);
}
export function getSolutionList(
  page,
  perpage,
  queryParam,
  username,
  language,
  result,
  contestId
) {
  return request('GET', 'solutions', {
    page: page,
    perpage: perpage,
    param: queryParam,
    username: username,
    language: language,
    result: result,
    contest_id: contestId
  });
}
export function getContestList(page, perpage, queryParam) {
  return request('GET', 'contests', {
    page: page,
    perpage: perpage,
    param: queryParam
  });
}
export function getContest(id) {
  return request('GET', `contest/${id}`);
}
export function getContestRankList(id) {
  return request('GET', `contest/${id}/ranklist`);
}
export function getContestTeamRankList(id) {
  return request('GET', `contest/${id}/teamranklist`);
}
export function getContestProblem(id, num) {
  return request('GET', `contest/${id}/problem/${num}`);
}
export function getSeriesList(page, perpage, queryParam) {
  return request('GET', 'serieses', {
    page: page,
    perpage: perpage,
    param: queryParam
  });
}
export function getSeries(id) {
  return request('GET', `series/${id}`);
}
export function getAllTags() {
  return request('GET', 'alltags');
}
