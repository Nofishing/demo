package demo

import (
	"awesomeProject/gin-vue-admin/server/model/system"
	"awesomeProject/pkg/mod/go.uber.org/zap@v1.16.0"
	"server/model/tpms"
)

func main()  {
	Excute()
}

func  Excute() (err error) {
	var userAuths []system.SysUseAuthority
	tx := global.GVA_DB
	//查找sys_user_auth里每个企业的所有角色
	err = tx.Model(&system.SysUseAuthority{}).
		Find(&userAuths).
		Group("enterprise_id, sys_authority_authority_id").
		Order("enterprise_id").Error
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		return err
	}

	var sysAuthorityMenusEnterprises []tpms.SysAuthorityMenusEnterprise
	//遍历按企业插入
	for _, v := range userAuths {
		//step1 先查模板select sys_authority_menus 这个地方可以优化成批量提交
		var sysAuthorityMenus []system.SysAuthorityMenus
		err = tx.Model(&system.SysAuthorityMenus{}).Where("sys_authority_authority_id = ?", v.SysAuthorityAuthorityId).Find(&sysAuthorityMenus).Error
		if err != nil {
			global.GVA_LOG.Error("获取失败!", zap.Error(err))
			return err
		}

		for _, va := range sysAuthorityMenus {
			var sysAuthorityMenusEnterprise tpms.SysAuthorityMenusEnterprise
			sysAuthorityMenusEnterprise.SysAuthorityAuthorityId = va.SysAuthorityAuthorityId
			sysAuthorityMenusEnterprise.SysBaseMenuId = va.SysBaseMenuId
			sysAuthorityMenusEnterprise.EnterpriseId = v.EnterpriseId
			sysAuthorityMenusEnterprises = append(sysAuthorityMenusEnterprises, sysAuthorityMenusEnterprise)
		}

	}

	//step2 再插入模板
	//插入新的用户相关菜单数据
	err = tx.Create(&sysAuthorityMenusEnterprises).Error
	return err
}
