## DB config

DB config is config of database. 

+ dialect: it should be `sqlite3` or `mysql`, only sqlite and mysql are supported for now.
+ db_path: db file path or mysql db config, eg(`root:12345678@(127.0.0.1:3306)/relayer?charset=utf8&parseTime=True&loc=Local`).

## Alert config

Relayer will send alert messages to telegram group if block is not be fetched for a long time or tx sent is failed.

+ moniker: `moniker` is moniker for relayer.
+ telegram_bot_id: `telegram_bot_id` is your telegram bot id.
+ telegram_chat_id: `telegram_chat_id` is chat id of group your bot joined.
+ block_update_time_out: `axc_block_update_time_out` is how long(in seconds) that block is not be fetched in asc chain you want 
relayer to send alert messages.

References:
+ [create a bot](https://core.telegram.org/bots#6-botfather)
+ [get bot id and chat id](https://stackoverflow.com/questions/32423837/telegram-bot-how-to-get-a-group-chat-id)

## Chain config

Chain common config for deputy. Pls note that `swap_amount` and `fixed_fee` below are number with decimal. For example, decimal in axim chain 
is 8 which means 100000000 is 1 actually. You need to handle decimal and amount with decimal.

+ asc_start_height: height of asc chain you want to start syncing when you start your relayer.
+ asc_providers: array of provider address of asc chain.
+ asc_confirm_num: confirm number of asc chain.
+ asc_token_hub_contract_address: token hub contract address of asc.
+ asc_validator_set_contract_address: validator set contract address of asc.

+ afc_rpc_addrs: array of rpc address of afc.
+ afc_key_type:  `mnemonic` and `aws_mnemonic` supported. `mnemonic` will use mnemonic provided below and `aws_mnemonic`
 will fetch mnemonic from aws secret manager.
+ afc_aws_region: region of aws.
+ afc_aws_secret_name: secret name of private key in aws.
+ afc_mnemonic: mnemonic of relayer operator.

## Log config

+ level: level of log, `CRITICAL`,`ERROR`,`WARNING`,`NOTICE`,`INFO`,`DEBUG` are supported.
+ filename: log file path if `use_console_logger` is true
+ max_file_size_in_mb: max log file size
+ max_backups_of_log_files: max backups of log files
+ max_age_to_retain_log_files_in_days: max days to retain log files
+ use_console_logger: use console logger or not
+ use_file_logger: use file logger or not
+ compress: compress log file or not
