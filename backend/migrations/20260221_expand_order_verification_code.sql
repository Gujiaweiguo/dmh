ALTER TABLE `orders`
MODIFY COLUMN `verification_code` VARCHAR(128) NULL COMMENT '核销码（包含签名）';
