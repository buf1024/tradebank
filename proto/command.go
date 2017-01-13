package proto

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

const (
	// 注意，所有请求报文都时奇数，所有应答报文都是偶数
	//////////////////////////////////////////////////////////////////////
	// 心跳类
	//////////////////////////////////////////////////////////////////////
	CMD_HEARTBEAT_REQ = 0x00010001 // 心跳请求
	CMD_HEARTBEAT_RSP = 0x00010002 // 心跳应答

	//////////////////////////////////////////////////////////////////////
	// 银行服务与资金管理模块交互
	//////////////////////////////////////////////////////////////////////
	CMD_E2E_CHANGE_ACCT_MONEY_REQ       = 0x07010001 // 资金变更请求
	CMD_E2E_CHANGE_ACCT_MONEY_RSP       = 0x07010002 // 资金变更应答
	CMD_E2E_CHANGE_INTERESTSET_REQ      = 0x07010003 // 交易所利息出入金请求
	CMD_E2E_CHANGE_INTERESTSET_RSP      = 0x07010004 // 交易所利息出入金应答
	CMD_E2E_CHANGE_ONE_SIDE_ACCOUNT_REQ = 0x07010005 // 单边账调整请求
	CMD_E2E_CHANGE_ONE_SIDE_ACCOUNT_RSP = 0x07010006 // 单边账调整应答
	CMD_E2E_CHANGE_COMMISION_REQ        = 0x07010009 // 手续费结转请求
	CMD_E2E_CHANGE_COMMISION_RSP        = 0x0701000A // 手续费结转应答

	//////////////////////////////////////////////////////////////////////
	// 银行服务与浮动盈亏模块交互
	//////////////////////////////////////////////////////////////////////
	CMD_E2E_QUERY_FLOAT_REVENUES_REQ = 0x22334487 // 查浮动盈亏请求
	CMD_E2E_QUERY_FLOAT_REVENUES_RSP = 0X22334488 // 查浮动盈亏应答

	//////////////////////////////////////////////////////////////////////
	// 银行服务与接口服务交互
	//////////////////////////////////////////////////////////////////////
	CMD_E2E_INOUT_MONEY_REQ          = 0x0F010001 // 出入金请求
	CMD_E2E_INOUT_MONEY_RSP          = 0x0F010002 // 出入金应答
	CMD_E2E_ATTACH_ACCT_REQ          = 0x0F010003 // 签约请求
	CMD_E2E_ATTACH_ACCT_RSP          = 0x0F010004 // 签约应答
	CMD_E2E_DETACH_ACCT_REQ          = 0x0F010005 // 解约请求
	CMD_E2E_DETACH_ACCT_RSP          = 0x0F010006 // 解约应答
	CMD_E2E_UPDATE_USER_INFO_REQ     = 0x0F010007 // 更新用户资料请求
	CMD_E2E_UPDATE_USER_INFO_RSP     = 0x0F010008 // 更新用户资料应答
	CMD_E2E_QUERY_MONEY_REQ          = 0x0F010009 // 查询交易所余额请求
	CMD_E2E_QUERY_MONEY_RSP          = 0x0F01000A // 查询交易所余额应答
	CMD_E2E_QUERY_SIGN_STATUS_REQ    = 0x0F01000B // 查询客户在银行签约状请求
	CMD_E2E_QUERY_SIGN_STATUS_RSP    = 0x0F01000C // 查询客户在银行签约状应答
	CMD_E2E_FINAL_FEE_REQ            = 0x0F01000D // 手续费结转请求
	CMD_E2E_FINAL_FEE_RSP            = 0x0F01000E // 手续费结转应答
	CMD_E2E_BRUTE_DEATTAH_REQ        = 0x0F01000F // 强解/强签请求
	CMD_E2E_BRUTE_DEATTAH_RSP        = 0x0F010010 // 强解/强签应答
	CMD_E2E_ONE_SIDE_ACCT_ADJUST_REQ = 0x0F010011 // 单边账调整请求
	CMD_E2E_ONE_SIDE_ACCT_ADJUST_RSP = 0x0F010012 // 单边账调整应答
	CMD_E2E_ONE_SIDE_ACCT_AUDIT_REQ  = 0x0F010013 // 单边账审批请求
	CMD_E2E_ONE_SIDE_ACCT_AUDIT_RSP  = 0x0F010014 // 单边账审批应答
	CMD_E2E_COMMISSION_CARRYOVER_REQ = 0x0F010015 // 手续费结转请求
	CMD_E2E_COMMISSION_CARRYOVER_RSP = 0x0F010016 // 手续费结转应答
	CMD_E2E_INTEREST_SETTLEMENT_REQ  = 0x0F010017 // 结息请求
	CMD_E2E_INTEREST_SETTLEMENT_RSP  = 0x0F010018 // 结息应答
	CMD_E2E_SIGNRESULT_NOTIFY_REQ    = 0x0F010019 // 签约结果通知请求
	CMD_E2E_SIGNRESULT_NOTIFY_RSP    = 0x0F01001A // 签约结果通知应答
	CMD_E2E_PAY_FORWARD_REQ          = 0x0F01001B // 支付推进请求
	CMD_E2E_PAY_FORWARD_RSP          = 0x0F01001C // 支付推进应答

	CMD_E2E_CLEAR_PROCESS_LINK_QUERY_REQ    = 0x0F0E0001 //清算处理环节查询请求
	CMD_E2E_CLEAR_PROCESS_LINK_QUERY_RSP    = 0x0F0E0002 //清算处理环节查询应答
	CMD_E2E_CLEAR_PROCESS_STATUS_QUERY_REQ  = 0x0F0E0003 //清算处理环节配置查询请求
	CMD_E2E_CLEAR_PROCESS_STATUS_QUERY_RSP  = 0x0F0E0004 //清算处理环节配置查询应答
	CMD_E2E_CLEAR_REPORT_SETTLE_QUERY_REQ   = 0x0F0E0005 //清算报表结算会员数据查询请求
	CMD_E2E_CLEAR_REPORT_SETTLE_QUERY_RSP   = 0x0F0E0006 //清算报表结算会员数据查询应答
	CMD_E2E_CLEAR_REPORT_MULTI_QUERY_REQ    = 0x0F0E0007 //清算报表综合会员数据查询请求
	CMD_E2E_CLEAR_REPORT_MULTI_QUERY_RSP    = 0x0F0E0008 //清算报表综合会员数据查询应答
	CMD_E2E_CLEAR_REPORT_CUSTOMER_QUERY_REQ = 0x0F0E0009 //清算报表交易客户数据查询请求
	CMD_E2E_CLEAR_REPORT_CUSTOMER_QUERY_RSP = 0x0F0E000A //清算报表交易客户数据查询应答
	CMD_E2E_CLEAR_REPORT_BROKER_QUERY_REQ   = 0x0F0E000B //清算报表经纪会员数据查询请求
	CMD_E2E_CLEAR_REPORT_BROKER_QUERY_RSP   = 0x0F0E000C //清算报表经纪会员数据查询应答
	CMD_E2E_CLEAR_RESULT_BOC_QUERY_REQ      = 0x0F0E0021 //中国银行的清算结果查询请求
	CMD_E2E_CLEAR_RESULT_BOC_QUERY_RSP      = 0x0F0E0022 //中国银行的清算结果查询应答
	CMD_E2E_CLEAR_RESULT_CCB_QUERY_REQ      = 0x0F0E0023 //建设银行的清算结果查询请求
	CMD_E2E_CLEAR_RESULT_CCB_QUERY_RSP      = 0x0F0E0024 //建设银行的清算结果查询应答

	// TODO 以下命令码在需要没确定之前预留
	CMD_E2E_SETTLEMENT_REQ        = 0x0A010001 // 清算请求
	CMD_E2E_SETTLEMENT_RSP        = 0x0A010002 // 清算应答
	CMD_E2E_SETTLEMENT_RESULT_REQ = 0x0A020001 // 清算结果查询请求
	CMD_E2E_SETTLEMENT_RESULT_RSP = 0x0A020002 // 清算结果查询应答
	CMD_E2E_TRANSFER_REQ          = 0x07020001 // 转账请求
	CMD_E2E_TRANSFER_RSP          = 0x07020002 // 转账应答
	CMD_E2E_UPDATE_BANK_CONF_REQ  = 0x03040003 // 更新银行配置请求
	CMD_E2E_UPDATE_BANK_CONF_RSP  = 0x03040004 // 更新银行配置应答
	CMD_E2E_QUERY_FINAL_FEE_REQ   = 0x83040003 // 查手续费结转请求
	CMD_E2E_QUERY_FINAL_FEE_RSP   = 0x83040004 // 查手续费结转应答
	CMD_E2E_QUERY_BANK_MONEY_REQ  = 0x83040093 // 查询银行余额请求
	CMD_E2E_QUERY_BANK_MONEY_RSP  = 0x83040094 // 查询银行余额应答
	// TODO 以上命令码在需要没确定之前预留

	//////////////////////////////////////////////////////////////////////
	// 交易系统与出入金服务交互
	//////////////////////////////////////////////////////////////////////
	CMD_E2B_SIGNINOUT_REQ                    = 0x11010001 // 银行服务向银行签到/签退请求
	CMD_E2B_SIGNINOUT_RSP                    = 0x11010002 // 银行服务向银行签到/签退应答
	CMD_E2B_ATTACH_ACCT_REQ                  = 0x11010003 // 银行服务向银行签约请求
	CMD_E2B_ATTACH_ACCT_RSP                  = 0x11010004 // 银行服务向银行签约应答
	CMD_E2B_DETACH_ACCT_REQ                  = 0x11010005 // 银行服务向银行解约请求
	CMD_E2B_DETACH_ACCT_RSP                  = 0x11010006 // 银行服务向银行解约应答
	CMD_E2B_IN_MONEY_REQ                     = 0x11010007 // 银行服务向银行入金请求
	CMD_E2B_IN_MONEY_RSP                     = 0x11010008 // 银行服务向银行入金应答
	CMD_E2B_OUT_MONEY_REQ                    = 0x11010009 // 银行服务向银行出金请求
	CMD_E2B_OUT_MONEY_RSP                    = 0x1101000A // 银行服务向银行出金应答
	CMD_E2B_QUERY_BANK_MONEY_REQ             = 0x1101000B // 银行服务向银行查询余额请求
	CMD_E2B_QUERY_BANK_MONEY_RSP             = 0x1101000C // 银行服务向银行查询余额入金应答
	CMD_E2B_ADJUST_MONEY_REQ                 = 0x1101000D // 银行服务向银行冲正请求
	CMD_E2B_ADJUST_MONEY_RSP                 = 0x1101000E // 银行服务向银行冲正应答
	CMD_E2B_FILE_NOTIFICATION_REQ            = 0x1101000F // 银行服务向银行文件通知请求
	CMD_E2B_FILE_NOTIFICATION_RSP            = 0x11010010 // 银行服务向银行文件通知应答
	CMD_E2B_UPDATE_USER_INFO_REQ             = 0x11010011 // 银行服务向银行更新用户资料请求
	CMD_E2B_UPDATE_USER_INFO_RSP             = 0x11010012 // 银行服务向银行更新用户资料应答
	CMD_E2B_QUERY_SIGN_STATUS_REQ            = 0x11010013 // 银行服务向银行客户在银行签约状请求
	CMD_E2B_QUERY_SIGN_STATUS_RSP            = 0x11010014 // 银行服务向银行客户在银行签约状应答
	CMD_E2B_OUT_MONEY_APPLICATION_RESULT_REQ = 0x11010015 // 出金审批申请结果请求 // 建行银行特殊流程
	CMD_E2B_OUT_MONEY_APPLICATION_RESULT_RSP = 0x11010016 // 出金审批申请结果应答 // 建行银行特殊流程
	CMD_E2B_CHECK_START_REQ                  = 0x11010017 // 银行服务向银行对账开始请求
	CMD_E2B_CHECK_START_RSP                  = 0x11010018 // 银行服务向银行对账开始应答
	CMD_E2B_CLEAR_REQ                        = 0x11010019 // 银行服务向银行清算开始请求
	CMD_E2B_CLEAR_RSP                        = 0x11010020 // 银行服务向银行清算开始应答
	CMD_E2B_PAY_FORWARD_REQ                  = 0x1101001B // 支付推进请求
	CMD_E2B_PAY_FORWARD_RSP                  = 0x1101001C // 支付推进应答

	//////////////////////////////////////////////////////////////////////
	// 出入金服务与交易系统交互
	//////////////////////////////////////////////////////////////////////
	CMD_B2E_ATTACH_ACCT_REQ         = 0x11020001 // 银行向银行服务签约请求
	CMD_B2E_ATTACH_ACCT_RSP         = 0x11020002 // 银行向银行服务签约应答
	CMD_B2E_DETACH_ACCT_REQ         = 0x11020003 // 银行向银行服务解约请求
	CMD_B2E_DETACH_ACCT_RSP         = 0x11020004 // 银行向银行服务解约应答
	CMD_B2E_IN_MONEY_REQ            = 0x11020005 // 银行向银行服务入金请求
	CMD_B2E_IN_MONEY_RSP            = 0x11020006 // 银行向银行服务入金应答
	CMD_B2E_OUT_MONEY_REQ           = 0x11020007 // 银行向银行服务出金请求
	CMD_B2E_OUT_MONEY_RSP           = 0x11020008 // 银行向银行服务出金应答
	CMD_B2E_QUERY_MONEY_REQ         = 0x11020009 // 银行向银行服务查询余额请求
	CMD_B2E_QUERY_MONEY_RSP         = 0x1102000A // 银行向银行服务查询余额应答
	CMD_B2E_ADJUST_MONEY_REQ        = 0x1102000B // 银行向银行冲正请求
	CMD_B2E_ADJUST_MONEY_RSP        = 0x1102000C // 银行向银行冲正应答
	CMD_B2E_FILE_NOTIFICATION_REQ   = 0x1102000D // 银行向银行文件通知请求
	CMD_B2E_FILE_NOTIFICATION_RSP   = 0x1102000E // 银行向银行文件通知应答
	CMD_B2E_UPDATE_USER_INFO_REQ    = 0x1102000F // 银行向银行变更用户属性请求
	CMD_B2E_UPDATE_USER_INFO_RSP    = 0x11020010 // 银行向银行变更用户属性应答
	CMD_B2E_QUERY_USER_PASSWORD_REQ = 0x11020011 // 银行端查询资金密码请求
	CMD_B2E_QUERY_USER_PASSWORD_RSP = 0x11020012 // 银行端查询资金密码应答

	/////////////////////////////////////////////////////////////////////
	// 建行银行特殊流程
	/////////////////////////////////////////////////////////////////////
	CMD_B2E_OUT_MONEY_APPLICATION_REQ = 0x11020013 // 出金审批申请请求
	CMD_B2E_OUT_MONEY_APPLICATION_RSP = 0x11020014 // 出金审批申请应答
	CMD_B2E_QUERY_OUT_MONEY_SID_REQ   = 0x11020015 // 会员出金申请流水查询请求
	CMD_B2E_QUERY_OUT_MONEY_SID_RSP   = 0x11020016 // 会员出金申请流水查询应答
	CMD_B2E_PUSH_USER_INFO_REQ        = 0x11020017 // 会员信息推送请求
	CMD_B2E_PUSH_USER_INFO_RSP        = 0x11020018 // 会员信息推送应答
	CMD_B2E_PUSH_INOUT_MONEY_INFO_REQ = 0x11020019 // 出入金推送请求
	CMD_B2E_PUSH_INOUT_MONEY_INFO_RSP = 0x1102001A // 出入金推送应答
	CMD_B2E_QUERY_USER_SIGN_INFO_REQ  = 0x1102001B // 会员签约信息查询请求
	CMD_B2E_QUERY_USER_SIGN_INFO_RSP  = 0x1102001C // 会员签约信息查询应答
	//////////////////////////////////////////////////////////////////////

	CMD_B2E_CHECK_FILE_NOTIFICATION_REQ = 0x1102001D // 银行端对账文件获取结果请求
	CMD_B2E_CHECK_FILE_NOTIFICATION_RSP = 0x1102001E // 银行端对账文件获取结果应答

	CMD_B2E_CLEAR_RESULT_REQ = 0x1102001F // 银行向银行服务通知清算结果请求
	CMD_B2E_CLEAR_RESULT_RSP = 0x11020020 // 银行向银行服务通知清算结果应答

	CMD_B2E_SIGN_INFO_REQ         = 0x11020021 // 银行向银行服务签约信息维护请求
	CMD_B2E_SIGN_INFO_RSP         = 0x11020022 // 银行向银行服务签约信息维护应答
	CMD_B2E_SUBACCOUNT_ATTACH_REQ = 0x11020023 // 银行向银行服务子帐户签约请求
	CMD_B2E_SUBACCOUNT_ATTACH_RSP = 0x11020024 // 银行向银行服务子帐户签约应答
	CMD_B2E_INOUTNOTIFY_REQ       = 0x11020025 // 银行端出入金推送  -- 网易宝请求
	CMD_B2E_INOUTNOTIFY_RSP       = 0x11020026 // 银行端出入金推送  -- 网易宝应答

	//////////////////////////////////////////////////////////////////////
	// 与统一接入服务平台服务注册
	//////////////////////////////////////////////////////////////////////
	CMD_SVR_REG_REQ = 0x0E000001 // 服务注册请求
	CMD_SVR_REG_RSP = 0x0E000002 // 服务注册应答
)

var message map[int64]proto.Message

func Message(command int64) (proto.Message, error) {
	if m, ok := message[command]; ok {
		proto.Clone(m)
		return m, nil
	}
	return nil, fmt.Errorf("mommand %d not found", command)
}

func Parse(command int64, buf []byte) (proto.Message, error) {
	msg, err := Message(command)
	if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(buf, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func Serialize(msg proto.Message) ([]byte, error) {
	return proto.Marshal(msg)
}

func Debug(command int64, msg proto.Message) string {
	return fmt.Sprintf("command:0x%x %s", command, proto.CompactTextString(msg))
}

func init() {
	message = make(map[int64]proto.Message)

	message[CMD_HEARTBEAT_REQ] = &HeartBeatReq{} //0x00010001 // 心跳请求
	message[CMD_HEARTBEAT_RSP] = &HeartBeatRsp{} //0x00010002                       // 心跳应答

	/*

		//////////////////////////////////////////////////////////////////////
		// 银行服务与资金管理模块交互
		//////////////////////////////////////////////////////////////////////
		message[CMD_E2E_CHANGE_ACCT_MONEY_REQ] =  &E2EChangeAcctMoneyReq{} //0x07010001       // 资金变更请求
		message[CMD_E2E_CHANGE_ACCT_MONEY_RSP] = &E2EChangeAcctMoneyRsp{} //0x07010002       // 资金变更应答
		message[CMD_E2E_CHANGE_INTERESTSET_REQ] = & //0x07010003      // 交易所利息出入金请求
		message[CMD_E2E_CHANGE_INTERESTSET_RSP] = & //0x07010004      // 交易所利息出入金应答
		message[CMD_E2E_CHANGE_ONE_SIDE_ACCOUNT_REQ] = & //0x07010005 // 单边账调整请求
		message[CMD_E2E_CHANGE_ONE_SIDE_ACCOUNT_RSP] = & //0x07010006 // 单边账调整应答
		message[CMD_E2E_CHANGE_COMMISION_REQ] = & //0x07010009        // 手续费结转请求
		message[CMD_E2E_CHANGE_COMMISION_RSP] = & //0x0701000A        // 手续费结转应答

		//////////////////////////////////////////////////////////////////////
		// 银行服务与浮动盈亏模块交互
		//////////////////////////////////////////////////////////////////////
		message[CMD_E2E_QUERY_FLOAT_REVENUES_REQ] = & //0x22334487 // 查浮动盈亏请求
		message[CMD_E2E_QUERY_FLOAT_REVENUES_RSP] = & //0x22334488 // 查浮动盈亏应答

		//////////////////////////////////////////////////////////////////////
		// 银行服务与接口服务交互
		//////////////////////////////////////////////////////////////////////
		message[CMD_E2E_INOUT_MONEY_REQ] = & //0x0F010001          // 出入金请求
		message[CMD_E2E_INOUT_MONEY_RSP] = & //0x0F010002          // 出入金应答
		message[CMD_E2E_ATTACH_ACCT_REQ] = & //0x0F010003          // 签约请求
		message[CMD_E2E_ATTACH_ACCT_RSP] = & //0x0F010004          // 签约应答
		message[CMD_E2E_DETACH_ACCT_REQ] = & //0x0F010005          // 解约请求
		message[CMD_E2E_DETACH_ACCT_RSP] = & //0x0F010006          // 解约应答
		message[CMD_E2E_UPDATE_USER_INFO_REQ] = & //0x0F010007     // 更新用户资料请求
		message[CMD_E2E_UPDATE_USER_INFO_RSP] = & //0x0F010008     // 更新用户资料应答
		message[CMD_E2E_QUERY_MONEY_REQ] = & //0x0F010009          // 查询交易所余额请求
		message[CMD_E2E_QUERY_MONEY_RSP] = & //0x0F01000A          // 查询交易所余额应答
		message[CMD_E2E_QUERY_SIGN_STATUS_REQ] = & //0x0F01000B    // 查询客户在银行签约状请求
		message[CMD_E2E_QUERY_SIGN_STATUS_RSP] = & //0x0F01000C    // 查询客户在银行签约状应答
		message[CMD_E2E_FINAL_FEE_REQ] = & //0x0F01000D            // 手续费结转请求
		message[CMD_E2E_FINAL_FEE_RSP] = & //0x0F01000E            // 手续费结转应答
		message[CMD_E2E_BRUTE_DEATTAH_REQ] = & //0x0F01000F        // 强解/强签请求
		message[CMD_E2E_BRUTE_DEATTAH_RSP] = & //0x0F010010        // 强解/强签应答
		message[CMD_E2E_ONE_SIDE_ACCT_ADJUST_REQ] = & //0x0F010011 // 单边账调整请求
		message[CMD_E2E_ONE_SIDE_ACCT_ADJUST_RSP] = & //0x0F010012 // 单边账调整应答
		message[CMD_E2E_ONE_SIDE_ACCT_AUDIT_REQ] = & //0x0F010013  // 单边账审批请求
		message[CMD_E2E_ONE_SIDE_ACCT_AUDIT_RSP] = & //0x0F010014  // 单边账审批应答
		message[CMD_E2E_COMMISSION_CARRYOVER_REQ] = & //0x0F010015 // 手续费结转请求
		message[CMD_E2E_COMMISSION_CARRYOVER_RSP] = & //0x0F010016 // 手续费结转应答
		message[CMD_E2E_INTEREST_SETTLEMENT_REQ] = & //0x0F010017  // 结息请求
		message[CMD_E2E_INTEREST_SETTLEMENT_RSP] = & //0x0F010018  // 结息应答
		message[CMD_E2E_SIGNRESULT_NOTIFY_REQ] = & //0x0F010019    // 签约结果通知请求
		message[CMD_E2E_SIGNRESULT_NOTIFY_RSP] = & //0x0F01001A    // 签约结果通知应答
		message[CMD_E2E_PAY_FORWARD_REQ] = & //0x0F01001B          // 支付推进请求
		message[CMD_E2E_PAY_FORWARD_RSP] = & //0x0F01001C          // 支付推进应答

		message[CMD_E2E_CLEAR_PROCESS_LINK_QUERY_REQ] = & //0x0F0E0001    //清算处理环节查询请求
		message[CMD_E2E_CLEAR_PROCESS_LINK_QUERY_RSP] = & //0x0F0E0002    //清算处理环节查询应答
		message[CMD_E2E_CLEAR_PROCESS_STATUS_QUERY_REQ] = & //0x0F0E0003  //清算处理环节配置查询请求
		message[CMD_E2E_CLEAR_PROCESS_STATUS_QUERY_RSP] = & //0x0F0E0004  //清算处理环节配置查询应答
		message[CMD_E2E_CLEAR_REPORT_SETTLE_QUERY_REQ] = & //0x0F0E0005   //清算报表结算会员数据查询请求
		message[CMD_E2E_CLEAR_REPORT_SETTLE_QUERY_RSP] = & //0x0F0E0006   //清算报表结算会员数据查询应答
		message[CMD_E2E_CLEAR_REPORT_MULTI_QUERY_REQ] = & //0x0F0E0007    //清算报表综合会员数据查询请求
		message[CMD_E2E_CLEAR_REPORT_MULTI_QUERY_RSP] = & //0x0F0E0008    //清算报表综合会员数据查询应答
		message[CMD_E2E_CLEAR_REPORT_CUSTOMER_QUERY_REQ] = & //0x0F0E0009 //清算报表交易客户数据查询请求
		message[CMD_E2E_CLEAR_REPORT_CUSTOMER_QUERY_RSP] = & //0x0F0E000A //清算报表交易客户数据查询应答
		message[CMD_E2E_CLEAR_REPORT_BROKER_QUERY_REQ] = & //0x0F0E000B   //清算报表经纪会员数据查询请求
		message[CMD_E2E_CLEAR_REPORT_BROKER_QUERY_RSP] = & //0x0F0E000C   //清算报表经纪会员数据查询应答
		message[CMD_E2E_CLEAR_RESULT_BOC_QUERY_REQ] = & //0x0F0E0021      //中国银行的清算结果查询请求
		message[CMD_E2E_CLEAR_RESULT_BOC_QUERY_RSP] = & //0x0F0E0022      //中国银行的清算结果查询应答
		message[CMD_E2E_CLEAR_RESULT_CCB_QUERY_REQ] = & //0x0F0E0023      //建设银行的清算结果查询请求
		message[CMD_E2E_CLEAR_RESULT_CCB_QUERY_RSP] = & //0x0F0E0024      //建设银行的清算结果查询应答

		// TODO 以下命令码在需要没确定之前预留
		message[CMD_E2E_SETTLEMENT_REQ] = & //0x0A010001        // 清算请求
		message[CMD_E2E_SETTLEMENT_RSP] = & //0x0A010002        // 清算应答
		message[CMD_E2E_SETTLEMENT_RESULT_REQ] = & //0x0A020001 // 清算结果查询请求
		message[CMD_E2E_SETTLEMENT_RESULT_RSP] = & //0x0A020002 // 清算结果查询应答
		message[CMD_E2E_TRANSFER_REQ] = & //0x07020001          // 转账请求
		message[CMD_E2E_TRANSFER_RSP] = & //0x07020002          // 转账应答
		message[CMD_E2E_UPDATE_BANK_CONF_REQ] = & //0x03040003  // 更新银行配置请求
		message[CMD_E2E_UPDATE_BANK_CONF_RSP] = & //0x03040004  // 更新银行配置应答
		message[CMD_E2E_QUERY_FINAL_FEE_REQ] = & //0x83040003   // 查手续费结转请求
		message[CMD_E2E_QUERY_FINAL_FEE_RSP] = & //0x83040004   // 查手续费结转应答
		message[CMD_E2E_QUERY_BANK_MONEY_REQ] = & //0x83040093  // 查询银行余额请求
		message[CMD_E2E_QUERY_BANK_MONEY_RSP] = & //0x83040094  // 查询银行余额应答
		// TODO 以上命令码在需要没确定之前预留

		//////////////////////////////////////////////////////////////////////
		// 交易系统与出入金服务交互
		//////////////////////////////////////////////////////////////////////
		message[CMD_E2B_SIGNINOUT_REQ] = & //0x11010001                    // 银行服务向银行签到/签退请求
		message[CMD_E2B_SIGNINOUT_RSP] = & //0x11010002                    // 银行服务向银行签到/签退应答
		message[CMD_E2B_ATTACH_ACCT_REQ] = & //0x11010003                  // 银行服务向银行签约请求
		message[CMD_E2B_ATTACH_ACCT_RSP] = & //0x11010004                  // 银行服务向银行签约应答
		message[CMD_E2B_DETACH_ACCT_REQ] = & //0x11010005                  // 银行服务向银行解约请求
		message[CMD_E2B_DETACH_ACCT_RSP] = & //0x11010006                  // 银行服务向银行解约应答
		message[CMD_E2B_IN_MONEY_REQ] = & //0x11010007                     // 银行服务向银行入金请求
		message[CMD_E2B_IN_MONEY_RSP] = & //0x11010008                     // 银行服务向银行入金应答
		message[CMD_E2B_OUT_MONEY_REQ] = & //0x11010009                    // 银行服务向银行出金请求
		message[CMD_E2B_OUT_MONEY_RSP] = & //0x1101000A                    // 银行服务向银行出金应答
		message[CMD_E2B_QUERY_BANK_MONEY_REQ] = & //0x1101000B             // 银行服务向银行查询余额请求
		message[CMD_E2B_QUERY_BANK_MONEY_RSP] = & //0x1101000C             // 银行服务向银行查询余额入金应答
		message[CMD_E2B_ADJUST_MONEY_REQ] = & //0x1101000D                 // 银行服务向银行冲正请求
		message[CMD_E2B_ADJUST_MONEY_RSP] = & //0x1101000E                 // 银行服务向银行冲正应答
		message[CMD_E2B_FILE_NOTIFICATION_REQ] = & //0x1101000F            // 银行服务向银行文件通知请求
		message[CMD_E2B_FILE_NOTIFICATION_RSP] = & //0x11010010            // 银行服务向银行文件通知应答
		message[CMD_E2B_UPDATE_USER_INFO_REQ] = & //0x11010011             // 银行服务向银行更新用户资料请求
		message[CMD_E2B_UPDATE_USER_INFO_RSP] = & //0x11010012             // 银行服务向银行更新用户资料应答
		message[CMD_E2B_QUERY_SIGN_STATUS_REQ] = & //0x11010013            // 银行服务向银行客户在银行签约状请求
		message[CMD_E2B_QUERY_SIGN_STATUS_RSP] = & //0x11010014            // 银行服务向银行客户在银行签约状应答
		message[CMD_E2B_OUT_MONEY_APPLICATION_RESULT_REQ] = & //0x11010015 // 出金审批申请结果请求 // 建行银行特殊流程
		message[CMD_E2B_OUT_MONEY_APPLICATION_RESULT_RSP] = & //0x11010016 // 出金审批申请结果应答 // 建行银行特殊流程
		message[CMD_E2B_CHECK_START_REQ] = & //0x11010017                  // 银行服务向银行对账开始请求
		message[CMD_E2B_CHECK_START_RSP] = & //0x11010018                  // 银行服务向银行对账开始应答
		message[CMD_E2B_CLEAR_REQ] = & //0x11010019                        // 银行服务向银行清算开始请求
		message[CMD_E2B_CLEAR_RSP] = & //0x11010020                        // 银行服务向银行清算开始应答
		message[CMD_E2B_PAY_FORWARD_REQ] = & //0x1101001B                  // 支付推进请求
		message[CMD_E2B_PAY_FORWARD_RSP] = & //0x1101001C                  // 支付推进应答

		//////////////////////////////////////////////////////////////////////
		// 出入金服务与交易系统交互
		//////////////////////////////////////////////////////////////////////
		message[CMD_B2E_ATTACH_ACCT_REQ] = & //0x11020001         // 银行向银行服务签约请求
		message[CMD_B2E_ATTACH_ACCT_RSP] = & //0x11020002         // 银行向银行服务签约应答
		message[CMD_B2E_DETACH_ACCT_REQ] = & //0x11020003         // 银行向银行服务解约请求
		message[CMD_B2E_DETACH_ACCT_RSP] = & //0x11020004         // 银行向银行服务解约应答
		message[CMD_B2E_IN_MONEY_REQ] = & //0x11020005            // 银行向银行服务入金请求
		message[CMD_B2E_IN_MONEY_RSP] = & //0x11020006            // 银行向银行服务入金应答
		message[CMD_B2E_OUT_MONEY_REQ] = & //0x11020007           // 银行向银行服务出金请求
		message[CMD_B2E_OUT_MONEY_RSP] = & //0x11020008           // 银行向银行服务出金应答
		message[CMD_B2E_QUERY_MONEY_REQ] = & //0x11020009         // 银行向银行服务查询余额请求
		message[CMD_B2E_QUERY_MONEY_RSP] = & //0x1102000A         // 银行向银行服务查询余额应答
		message[CMD_B2E_ADJUST_MONEY_REQ] = & //0x1102000B        // 银行向银行冲正请求
		message[CMD_B2E_ADJUST_MONEY_RSP] = & //0x1102000C        // 银行向银行冲正应答
		message[CMD_B2E_FILE_NOTIFICATION_REQ] = & //0x1102000D   // 银行向银行文件通知请求
		message[CMD_B2E_FILE_NOTIFICATION_RSP] = & //0x1102000E   // 银行向银行文件通知应答
		message[CMD_B2E_UPDATE_USER_INFO_REQ] = & //0x1102000F    // 银行向银行变更用户属性请求
		message[CMD_B2E_UPDATE_USER_INFO_RSP] = & //0x11020010    // 银行向银行变更用户属性应答
		message[CMD_B2E_QUERY_USER_PASSWORD_REQ] = & //0x11020011 // 银行端查询资金密码请求
		message[CMD_B2E_QUERY_USER_PASSWORD_RSP] = & //0x11020012 // 银行端查询资金密码应答

		/////////////////////////////////////////////////////////////////////
		// 建行银行特殊流程
		/////////////////////////////////////////////////////////////////////
		message[CMD_B2E_OUT_MONEY_APPLICATION_REQ] = & //0x11020013 // 出金审批申请请求
		message[CMD_B2E_OUT_MONEY_APPLICATION_RSP] = & //0x11020014 // 出金审批申请应答
		message[CMD_B2E_QUERY_OUT_MONEY_SID_REQ] = & //0x11020015   // 会员出金申请流水查询请求
		message[CMD_B2E_QUERY_OUT_MONEY_SID_RSP] = & //0x11020016   // 会员出金申请流水查询应答
		message[CMD_B2E_PUSH_USER_INFO_REQ] = & //0x11020017        // 会员信息推送请求
		message[CMD_B2E_PUSH_USER_INFO_RSP] = & //0x11020018        // 会员信息推送应答
		message[CMD_B2E_PUSH_INOUT_MONEY_INFO_REQ] = & //0x11020019 // 出入金推送请求
		message[CMD_B2E_PUSH_INOUT_MONEY_INFO_RSP] = & //0x1102001A // 出入金推送应答
		message[CMD_B2E_QUERY_USER_SIGN_INFO_REQ] = & //0x1102001B  // 会员签约信息查询请求
		message[CMD_B2E_QUERY_USER_SIGN_INFO_RSP] = & //0x1102001C  // 会员签约信息查询应答
		//////////////////////////////////////////////////////////////////////

		message[CMD_B2E_CHECK_FILE_NOTIFICATION_REQ] = & //0x1102001D // 银行端对账文件获取结果请求
		message[CMD_B2E_CHECK_FILE_NOTIFICATION_RSP] = & //0x1102001E // 银行端对账文件获取结果应答

		message[CMD_B2E_CLEAR_RESULT_REQ] = & //0x1102001F // 银行向银行服务通知清算结果请求
		message[CMD_B2E_CLEAR_RESULT_RSP] = & //0x11020020 // 银行向银行服务通知清算结果应答

		message[CMD_B2E_SIGN_INFO_REQ] = & //0x11020021         // 银行向银行服务签约信息维护请求
		message[CMD_B2E_SIGN_INFO_RSP] = & //0x11020022         // 银行向银行服务签约信息维护应答
		message[CMD_B2E_SUBACCOUNT_ATTACH_REQ] = & //0x11020023 // 银行向银行服务子帐户签约请求
		message[CMD_B2E_SUBACCOUNT_ATTACH_RSP] = & //0x11020024 // 银行向银行服务子帐户签约应答
		message[CMD_B2E_INOUTNOTIFY_REQ] = & //0x11020025       // 银行端出入金推送  -- 网易宝请求
		message[CMD_B2E_INOUTNOTIFY_RSP] = & //0x11020026       // 银行端出入金推送  -- 网易宝应答
	*/
	//////////////////////////////////////////////////////////////////////
	// 与统一接入服务平台服务注册
	//////////////////////////////////////////////////////////////////////
	message[CMD_SVR_REG_REQ] = &SvrRegReq{} //0x0E000001 // 服务注册请求
	message[CMD_SVR_REG_RSP] = &SvrRegRsp{} //0x0E000002 // 服务注册应答

}
