--建库以及授权--
create database SCMPLAT;
grant all privileges on SCMPLAT.* to 'SCMPLAT'@'localhost' identified by 'ScmPlat-123';
grant all privileges on SCMPLAT.* to 'SCMPLAT'@'%' identified by 'ScmPlat-123';
grant all privileges on SCMPLAT.* to 'AIDBA'@'localhost' identified by 'AiDBA-123456';
grant all privileges on SCMPLAT.* to 'AIDBA'@'%' identified by 'AiDBA-123456';
flush privileges;

---权限初始化---
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('1', 'admin', 'SecMg', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('2', 'admin', 'ProcAddButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('3', 'admin', 'ProcModifyButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('4', 'admin', 'ProcDelButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('5', 'admin', 'ProcCreateButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('6', 'admin', 'PubAddButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('7', 'admin', 'PubModifyButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('8', 'admin', 'PubDelButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('9', 'admin', 'PubPublishButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('10', 'admin', 'StAddButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('11', 'admin', 'StModifyButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('12', 'admin', 'StDelButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('13', 'admin', 'StExecButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('14', 'operator', 'SecMg', 'none');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('15', 'operator', 'ProcAddButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('16', 'operator', 'ProcModifyButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('17', 'operator', 'ProcDelButton', 'none');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('18', 'operator', 'ProcCreateButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('19', 'operator', 'PubAddButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('20', 'operator', 'PubModifyButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('21', 'operator', 'PubDelButton', 'none');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('22', 'operator', 'PubPublishButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('23', 'operator', 'StAddButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('24', 'operator', 'StModifyButton', 'inline-block');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('25', 'operator', 'StDelButton', 'none');
INSERT INTO `SCMPLAT`.`sec_mg` (`sec_id`, `role_name`, `oper_id`, `dpl_status`) VALUES ('26', 'operator', 'StExecButton', 'inline-block');

