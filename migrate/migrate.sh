mysql -h db -uroot -p$MYSQL_PASSWORD <<EOF
truncate table oj.user;
insert into oj.$(user) (email,username,nick,avatar,passsalt,password,submit,solved,defunct,created_at,updated_at,is_compete_user,role_id)
select ou.email,ou.user_id,ou.user_id,'static/images/default_avatar.png','abcdefghijklmnop','72bd4f7ef86a8492eba63b96ee11d1f329fff8cc',ou.submit,ou.solved,0,ou.reg_time,NOW(),0,1 from jol.users ou where ou.user_id not like 'A%' and ou.user_id not like 'B%' and ou.user_id not like 'test%' and ou.user_id not like 'AHPU%' and ou.user_id not like 'AHUT%';

truncate table oj.problem;
insert into oj.problem (id,title,description,level,input,output,sample_input,sample_output,spj,hint,defunct,time_limit,memory_limit,accepted,submit,solved,created_at,updated_at)
select op.problem_id-1000,op.title,op.description,1,op.input,op.output,op.sample_input,op.sample_output,0,op.hint,0,op.time_limit,op.memory_limit,op.accepted,op.submit,op.solved,op.in_date,op.in_date from jol.problem op;

DROP TEMPORARY TABLE IF EXISTS results;
create temporary table results select u.id,u.username,s.* from jol.solution s inner join oj.user u on s.user_id = u.username  where s.result != 13;

DROP TEMPORARY TABLE IF EXISTS src;
create temporary table src
(
    id int primary key auto_increment,
    source text
)
select src.source from jol.solution s inner join jol.source_code src on s.solution_id = src.solution_id  inner join oj.user u on s.user_id = u.username where s.result != 13;

truncate table oj.solution;
insert into oj.solution (problem_id,team_id,user_id,contest_id,num,time,memory,in_date,result,language,ip,judgetime,valid,code_length,pass_rate,lint_error,judger)
select os.problem_id-1000,0,os.id,0,os.num,os.time,os.memory,os.in_date,os.result,os.language,os.ip,os.judgetime,1,os.code_length,os.pass_rate,os.lint_error,os.judger from results os;

DROP TEMPORARY TABLE IF EXISTS results;

truncate table oj.source_code;
insert into oj.source_code
select src.id,src.source,0 from src;

update user set email = null where email = '';
delete user from user left join (select email,min(id) as min_id from user group by email) u1 on user.id =u1.min_id where u1.min_id is null and user.email != '';

EOF
