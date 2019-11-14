CREATE TABLE capture
(
    id         INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    file_name  VARCHAR(32)        NOT NULL,
    created_at TIMESTAMP          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP          NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP          NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE shoe_box
(
    id         INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    is_exist   boolean            NOT NULL, -- 存在しているか
    live_times TINYINT            NOT NULL, -- 連続して存在していた回数
    row        TINYINT            NOT NULL, -- 靴箱の行番号
    `column`   TINYINT            NOT NULL, -- 靴箱の列番号
    capture_id INT                NOT NULL,
    created_at TIMESTAMP          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP          NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP          NULL,
    FOREIGN KEY (capture_id) REFERENCES capture (id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
