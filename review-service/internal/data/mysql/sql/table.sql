-- Active: 1746952367450@@127.0.0.1@13308@review
CREATE TABLE review_info (
    `id` bigint(32) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    `create_by` VARCHAR(48) NOT NULL DEFAULT '' COMMENT '创建方式标识',
    `update_by` VARCHAR(48) NOT NULL DEFAULT '' COMMENT '更新方式标识',
    `create_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_at` TIMESTAMP NULL DEFAULT NULL COMMENT '逻辑删除标记',
    `version` INT(10) NOT NULL DEFAULT 0 COMMENT '乐观锁标记',

    `review_id` bigint(32) UNSIGNED NOT NULL DEFAULT 0 COMMENT '评论id',
    `content` VARCHAR(512) NOT NULL COMMENT '评价内容',
    `socore` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '评分',
    `service_score` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '服务评分',
    `express_score` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '物流评分',
    `has_media` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '是否有图或视频',
    `order_id` bigint(32) UNSIGNED NOT NULL DEFAULT 0 COMMENT '订单id',
    `sku_id` bigint(32) UNSIGNED NOT NULL DEFAULT 0 COMMENT '商品id',
    `spu_id` bigint(32) UNSIGNED NOT NULL DEFAULT 0 COMMENT '货号',
    `store_id` bigint(32) UNSIGNED NOT NULL DEFAULT 0 COMMENT '店铺id',
    `user_id` bigint(32) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
    `anonymous` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '是否匿名',
    `tags` VARCHAR(1024) NOT NULL COMMENT '标签json',
    `pic_info` VARCHAR(1024) NOT NULL COMMENT '图片信息json',
    `video_info` VARCHAR(1024) NOT NULL COMMENT '视频信息json',
    `status` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '状态:10-待审核，20-审核通过，30-审核不通过，40-隐藏',
    `is_default` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '是否默认评价',
    `has_reply` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '是否有回复',
    `op_reason` VARCHAR(512) NOT NULL COMMENT '审核拒绝原因',
    `op_remark` VARCHAR(512) NOT NULL COMMENT '审核备注',
    `op_user` VARCHAR(48) NOT NULL DEFAULT '' COMMENT '审核人',
    `goods_snapshoot` VARCHAR(1024) NOT NULL COMMENT '商品快照',
    
    `ext_json` VARCHAR(1024) NOT NULL COMMENT '扩展信息json',
    `ctrl_json` VARCHAR(1024) NOT NULL COMMENT '控制信息json',

    PRIMARY KEY (`id`),
    KEY `idx_delete_at` (`delete_at`) COMMENT '逻辑删除标记索引',
    UNIQUE KEY `uk_review_id` (`review_id`) COMMENT '评论id唯一索引',
    KEY `idx_order_id` (`order_id`) COMMENT '订单id索引',
    key `idx_user_id` (`user_id`) COMMENT '用户id索引'
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评价信息表';

CREATE Table review_reply_info (
    `id` bigint(32) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    `create_by` VARCHAR(48) NOT NULL DEFAULT '' COMMENT '创建方式标识',
    `update_by` VARCHAR(48) NOT NULL DEFAULT '' COMMENT '更新方式标识',
    `create_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_at` TIMESTAMP NULL DEFAULT NULL COMMENT '逻辑删除标记',
    `version` INT(10) NOT NULL DEFAULT 0 COMMENT '乐观锁标记',

    `reply_id` bigint(32) UNSIGNED NOT NULL DEFAULT 0 COMMENT '回复id',
    `review_id` bigint(32) UNSIGNED NOT NULL DEFAULT 0 COMMENT '评论id',
    `store_id` bigint(32) UNSIGNED NOT NULL DEFAULT 0 COMMENT '店铺id',
    `content` VARCHAR(512) NOT NULL COMMENT '回复内容',
    `pic_info` VARCHAR(1024) NOT NULL COMMENT '图片信息json',
    `video_info` VARCHAR(1024) NOT NULL COMMENT '视频信息json',

    `ext_json` VARCHAR(1024) NOT NULL COMMENT '扩展信息json',
    `ctrl_json` VARCHAR(1024) NOT NULL COMMENT '控制信息json',

    PRIMARY KEY (`id`),
    KEY `idx_delete_at` (`delete_at`) COMMENT '逻辑删除标记索引',
    UNIQUE KEY `uk_reply_id` (`reply_id`) COMMENT '回复id唯一索引',
    KEY `idx_review_id` (`review_id`) COMMENT '评论id索引',
    KEY `idx_store_id` (`store_id`) COMMENT '店铺id索引'
)engine=InnoDB DEFAULT charset=utf8mb4 COMMENT='评价回复信息表';

CREATE TABLE review_appeal_info (
    `id` bigint(32) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    `create_by` VARCHAR(48) NOT NULL DEFAULT '' COMMENT '创建方式标识',
    `update_by` VARCHAR(48) NOT NULL DEFAULT '' COMMENT '更新方式标识',
    `create_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_at` TIMESTAMP NULL DEFAULT NULL COMMENT '逻辑删除标记',
    `version` INT(10) NOT NULL DEFAULT 0 COMMENT '乐观锁标记',

    `appeal_id` bigint(32) UNSIGNED NOT NULL DEFAULT 0 COMMENT '申诉id',
    `review_id` bigint(32) UNSIGNED NOT NULL DEFAULT 0 COMMENT '评论id',
    `store_id` bigint(32) UNSIGNED NOT NULL DEFAULT 0 COMMENT '店铺id',
    `status` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '状态:10-待处理，20-处理中，30-已处理',
    `reason` VARCHAR(512) NOT NULL COMMENT '申诉原因',
    `content` VARCHAR(512) NOT NULL COMMENT '申诉内容',
    `pic_info` VARCHAR(1024) NOT NULL COMMENT '图片信息json',
    `video_info` VARCHAR(1024) NOT NULL COMMENT '视频信息json',

    `op_remark` VARCHAR(512) NOT NULL COMMENT '处理备注',
    `op_user` VARCHAR(48) NOT NULL DEFAULT '' COMMENT '处理人',
    
    `ext_json` VARCHAR(1024) NOT NULL COMMENT '扩展信息json',
    `ctrl_json` VARCHAR(1024) NOT NULL COMMENT '控制信息json',

    PRIMARY KEY (`id`),
    KEY `idx_delete_at` (`delete_at`) COMMENT '逻辑删除标记索引',
    KEY `idx_appeal_id` (`appeal_id`) COMMENT '申诉id索引',
    UNIQUE KEY `uk_review_id` (`review_id`) COMMENT '评论id唯一索引',
    KEY `idx_store_id` (`store_id`) COMMENT '店铺id索引'
)engine=InnoDB DEFAULT charset=utf8mb4 COMMENT='评价申诉信息表';