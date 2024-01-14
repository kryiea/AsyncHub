/*
    - MySQL的utf8是utfmb3，只有三个字节，节省空间但不能表达全部的UTF-8。所以推荐使用utf8mb4。
    - utf8mb4_unicode_ci,是基于标准的Unicode来排序和比较，能够在各种语言之间精确排序
    - ci 表示不区分大小写，也就是说，排序时 a 和 A 之间没有区别
 */

# 任务信息表
CREATE TABLE `t_lark_task_1`(
        `id` INT(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
        `user_id` VARCHAR(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '用户id',
        `task_id` VARCHAR(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '任务id',
        `task_type` VARCHAR(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '任务类型',
        `task_stage` VARCHAR(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '任务阶段',
        `status` TINYINT(3) UNSIGNED NOT NULL DEFAULT '1' COMMENT '任务状态',
        `priority` INT(11) NOT NULL DEFAULT '0' COMMENT '任务优先级',
        `crt_retry_num` INT(11) NOT NULL DEFAULT '0' COMMENT '当前重试次数',
        `max_retry_num` INT(11) NOT NULL DEFAULT '0' COMMENT '最大重试次数',
        `max_retry_interval` INT(11) NOT NULL DEFAULT '0' COMMENT '最大重试间隔',
        `schedule_log` VARCHAR(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '调度日志',
        `task_context` VARCHAR(8192) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '任务上下文',
        `order_time` INT(20) NOT NULL DEFAULT '0' COMMENT '调度时间，越小越优先',
        `create_time` datetime NULL DEFAULT current_timestamp COMMENT '创建时间',
        `modify_time` datetime NULL DEFAULT current_timestamp ON UPDATE current_timestamp COMMENT '修改时间',
        PRIMARY KEY(`id`),
        UNIQUE KEY `idx_task_id` (`task_id`),
        KEY `idx_user_id` (`user_id`),
        KEY `idx_tasktype_status_modify_time` (`task_type,` `status`, `modify_time`)
)ENGINE = InnoDB DEFAULt CHARSET = utf8mb4 COLLATE utf8mb4_unicode_ci;


# 任务调度配置表
# TODO：创建时间为什么可以为空？
CREATE TABLE `t_schedule_cfg`(
        `task_type` VARCHAR(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '任务类型',
        `schedule_limit` INT(11) DEFAULT '0' COMMENT '每次拉取的最大任务数',
        `schedule_interval` INT(11) DEFAULT '10' COMMENT '调度间隔',
        `max_processing_time` INT(11) DEFAULT '0' COMMENT '最大处理时间',
        `max_retry_num` INT(11) DEFAULT '0' COMMENT '最大重试次数',
        `retry_interval` INT(11) DEFAULT NULL COMMENT '初始重试间隔',
        `max_retry_interval` INT(11) DEFAULT NULL COMMENT '最大重试间隔',
        `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
        `modify_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
        PRIMARY KEY(`task_type`)
)ENGINE = InnoDB DEFAULt CHARSET = utf8mb4 COLLATE utf8mb4_unicode_ci;

# 任务调度位置表
CREATE TABLE `t_schedule_pos`(
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `task_type` VARCHAR(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '任务类型',
    `schedule_begin_pos` INT(11) NOT NULL DEFAULT '0' COMMENT '调度开始于几号表',
    `schedule_end_pos` INT(11) NOT NULL DEFAULT '0' COMMENT '调度结束于几号表',
    `create_time` datetime NOT NULL DEFAULT current_timestamp COMMENT '创建时间',
    `modify_time` datetime NOT NULL DEFAULT current_timestamp ON UPDATE current_timestamp COMMENT '修改时间',
    PRIMARY KEY(`id`),
    UNIQUE KEY `idx_task_type` (`task_type`)
)ENGINE = InnoDB DEFAULt CHARSET = utf8mb4 COLLATE utf8mb4_unicode_ci;

insert into t_schedule_cfg(
        task_type,
        schedule_limit,
        schedule_interval,
        max_processing_time,
        max_retry_num,
        retry_interval,
        max_retry_interval
    )
    values("lark", 100, 10, 30, 3, 5, 30);


insert into t_schedule_pos(
        task_type,
        schedule_begin_pos,
        schedule_end_pos
    )
    values("lark", 1, 1);