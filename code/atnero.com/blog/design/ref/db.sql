CREATE
    TABLE dbversion(
        version CHAR(8)
        ,CONSTRAINT pk_dbversion PRIMARY KEY (version)
    );

CREATE
    TABLE account(
        account_id AUTOINCREMENT UNIQUE
        ,mac CHAR(12) NOT null
        ,uuid CHAR(32) NOT null
        ,ip VARCHAR(40) not null default '0.0.0.0'
        ,login_name VARCHAR(40) WITH COMP
        ,login_pwd VARCHAR(40) WITH COMP
        ,reg_time DATETIME NOT null DEFAULT now()
        ,lastlogin_time DATETIME
        ,lastonline_time DATETIME
        ,online bit DEFAULT false NOT null
        ,CONSTRAINT pk_account PRIMARY KEY (
            mac
            ,uuid
        )
    );
    
CREATE
    TABLE winctrl_rtc(
        rtc_id autoincrement
        ,mode INTEGER NOT null
        ,CONSTRAINT chk_mode CHECK(mode IN (
            0 -- Data
            ,1 -- Day
            ,255 -- 0xff disable
        ))
        ,wu_date INTEGER NOT null -- Date
        ,constraint chk_wu_date CHECK(wu_date between 0 and 31)
        ,wu_week INTEGER NOT null
        ,constraint chk_wu_week CHECK(wu_week between 0 and 255)	--bit 1 - bit 7
        ,wu_hour INTEGER NOT null
        ,constraint chk_wu_hour CHECK(wu_hour between 0 and 23)
        ,wu_minute INTEGER NOT null
        ,constraint chk_wu_minute CHECK(wu_minute between 0 and 59)
        ,wu_second INTEGER NOT null
        ,constraint chk_wu_second CHECK(wu_second between 0 and 59)
        ,account_id LONG NOT null
        ,CONSTRAINT pk_winctrl PRIMARY KEY (rtc_id)
        ,CONSTRAINT fk_winctrl FOREIGN KEY (account_id) references account(account_id) on delete cascade
    );

CREATE
    TABLE winctrl_tpo(
        tpo_id autoincrement
        ,mode INTEGER NOT null
        ,CONSTRAINT tpo_chk_mode CHECK(mode IN (
            0
            ,1
            ,255
        ))
        ,wu_date INTEGER NOT null
        ,constraint tpo_chk_wu_date CHECK(wu_date between 0 and 31)
        ,wu_week INTEGER NOT null
        ,constraint tpo_chk_wu_week CHECK(wu_week between 0 and 255)
        ,wu_hour INTEGER NOT null
        ,constraint tpo_chk_wu_hour CHECK(wu_hour between 0 and 23)
        ,wu_minute INTEGER NOT null
        ,constraint tpo_chk_wu_minute CHECK(wu_minute between 0 and 59)
        ,wu_second INTEGER NOT null
        ,constraint tpo_chk_wu_second CHECK(wu_second between 0 and 59)
        ,account_id LONG NOT null
        ,CONSTRAINT pk_winctrl_tpo PRIMARY KEY (tpo_id)
        ,CONSTRAINT fk_winctrl_tpo FOREIGN KEY (account_id) references account(account_id) on delete cascade
    );

CREATE
    TABLE winctrl_rtc(
        rtc_id autoincrement
        ,mode INTEGER NOT null
        ,CONSTRAINT chk_mode CHECK(mode IN (
            0
            ,1
            ,255
        ))
        ,wu_date INTEGER NOT null
        ,constraint chk_wu_date CHECK(wu_date between 0 and 31)
        ,wu_week INTEGER NOT null
        ,constraint chk_wu_week CHECK(wu_week between 0 and 255)
        ,wu_hour INTEGER NOT null
        ,constraint chk_wu_hour CHECK(wu_hour between 0 and 23)
        ,wu_minute INTEGER NOT null
        ,constraint chk_wu_minute CHECK(wu_minute between 0 and 59)
        ,wu_second INTEGER NOT null
        ,constraint chk_wu_second CHECK(wu_second between 0 and 59)
        ,account_id LONG NOT null
        ,CONSTRAINT pk_winctrl PRIMARY KEY (rtc_id)
        ,CONSTRAINT fk_winctrl FOREIGN KEY (account_id) references account(account_id) on delete cascade
    );
    
alter table winctrl_rtc drop constraint chk_mode
alter table winctrl_rtc drop constraint chk_wu_date
alter table winctrl_rtc drop constraint chk_wu_week
alter table winctrl_rtc drop constraint chk_wu_hour
alter table winctrl_rtc drop constraint chk_wu_minute
alter table winctrl_rtc drop constraint chk_wu_second

alter table winctrl_rtc add constraint chk_mode check(mode in (0,1,255))

CREATE
    TABLE hwm(
        hwm_id AUTOINCREMENT
        ,becheck bit DEFAULT false NOT null
        ,account_id LONG NOT null
        ,CONSTRAINT pk_hwm PRIMARY KEY (hwm_id)
        ,CONSTRAINT fk_hwm FOREIGN KEY (account_id) references account(account_id) on delete cascade
    );
    
CREATE
    TABLE board_module(
        bm_id AUTOINCREMENT
        ,project_path varchar(10)
        ,board_id varchar(30)
        ,maj_ver smallint not null default 0
        ,min_ver smallint not null default 0
        ,account_id LONG NOT null
        ,CONSTRAINT pk_board_module PRIMARY KEY (bm_id)
        ,CONSTRAINT fk_board_module FOREIGN KEY (account_id) references account(account_id) on delete cascade
    );
    
    
INSERT
    INTO
        account(
            mac
            ,uuid
            ,login_name
        )
    VALUES
        (
            '112233445566'
            ,'0123456789ABCDEF0123456789ABCDEF'
            ,'XXXYYY'
        );
        
