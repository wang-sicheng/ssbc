DROP TABLE IF EXISTS `block`;
CREATE TABLE `block` (
  `id`  bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `pid` bigint(20) NOT NULL DEFAULT 0 COMMENT '父块id',
  `prev_hash`  varchar(255) NOT NULL DEFAULT 'prev_hash' COMMENT '父块哈希',
  `hash` varchar(255) NOT NULL DEFAULT 'hash' COMMENT '区块哈希',
  `merkle_root` varchar(255) NOT NULL DEFAULT 'merkle_root' COMMENT '区块交易的默克尔树根',
  `tx_count` bigint(20) NOT NULL DEFAULT 0 COMMENT '包含的交易数量',
  `signature` varchar(255) NOT NULL DEFAULT 'signature' COMMENT '打包者签名',
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '时间戳',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='区块表';

DROP TABLE IF EXISTS `transaction`;
CREATE TABLE `transaction` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `block_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '区块id',
  `sender_address` varchar(255) NOT NULL DEFAULT 'sender_address' COMMENT '发起者地址',
  `receiver_address` varchar(255) NOT NULL DEFAULT 'receiver_address' COMMENT '接收者地址',
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '时间戳',
  `signature` varchar(255) NOT NULL DEFAULT 'signature' COMMENT '发起者签名',
  `message` varchar(255) DEFAULT NULL DEFAULT 'message' COMMENT '消息（暂不知道什么用）',
  `sender_public_key` varchar(255) NOT NULL DEFAULT 'public_key' COMMENT '发起者公钥',
  `transfer_amount` bigint(10) NOT NULL DEFAULT 0 COMMENT '转账金额',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='交易表';

