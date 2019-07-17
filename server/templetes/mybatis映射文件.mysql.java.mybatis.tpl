<?xml version="1.0" encoding="UTF-8" ?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd" >
<mapper namespace="com.jwell.msa.customer.provider.account.dao.AccountMapper">
    <resultMap id="BaseResultMap" type="com.jwell.msa.customer.provider.account.entity.AccountDO">
        <id column="sysno" property="sysno" jdbcType="BIGINT"/>
        <result column="login_name" property="loginName" jdbcType="VARCHAR"/>
        <result column="mobile" property="mobile" jdbcType="VARCHAR"/>
        <result column="account_type" property="accountType" jdbcType="TINYINT"/>
        <result column="account_from_sys" property="accountFromSys" jdbcType="VARCHAR"/>
        <result column="account_status" property="accountStatus" jdbcType="TINYINT"/>
        <result column="invite_code" property="inviteCode" jdbcType="VARCHAR"/>
        <result column="invite_account_sysno" property="inviteAccountSysno" jdbcType="BIGINT"/>
        <result column="nick_name" property="nickName" jdbcType="VARCHAR"/>
        <result column="email" property="email" jdbcType="VARCHAR"/>
        <result column="email_verify" property="emailVerify" jdbcType="TINYINT"/>
        <result column="user_memcode" property="userMemCode" jdbcType="VARCHAR"/>
        <result column="create_date" property="createDate" jdbcType="TIMESTAMP"/>
        <result column="create_user" property="createUser" jdbcType="VARCHAR"/>
        <result column="modify_date" property="modifyDate" jdbcType="TIMESTAMP"/>
        <result column="modify_user" property="modifyUser" jdbcType="VARCHAR"/>
        <result column="administrator" property="administrator" jdbcType="TINYINT"/>
    </resultMap>

    <select id="selectByPrimaryKey" resultMap="BaseResultMap" parameterType="java.lang.Long">
    select 
        `sysno`,
        `login_name`,
        `mobile`,
        `account_type`,
        `account_from_sys`,
        `account_status`,
        `invite_code`,
        `invite_account_sysno`,
        `nick_name`,
        `email`,
        `email_verify`,
        `user_memcode`,
        `create_date`,
        `create_user`,
        `modify_date`,
        `modify_user`
    from
        `account`
    where
        `sysno` = #{sysno,jdbcType=BIGINT}
  </select>

    <insert id="insert" parameterType="com.jwell.msa.customer.provider.account.entity.AccountDO" keyProperty="sysno"
            useGeneratedKeys="true">
    insert into account (
        `sysno`,
        `login_name`,
        `mobile`,
        `account_type`,
        `account_from_sys`,
        `account_status`,
        `invite_code`,
        `invite_account_sysno`,
        `nick_name`,
        `email`,
        `email_verify`,
        `user_memcode`,
        `create_user`,
        `modify_date`,
        `modify_user`
    )
    values (
          #{sysno,jdbcType=BIGINT},
          #{loginName,jdbcType=VARCHAR},
          #{mobile,jdbcType=VARCHAR},
          #{accountType,jdbcType=TINYINT},
          #{accountFromSys,jdbcType=VARCHAR},
          #{accountStatus,jdbcType=TINYINT},
          #{inviteCode,jdbcType=VARCHAR},
          #{inviteAccountSysno,jdbcType=BIGINT},
          #{nickName,jdbcType=VARCHAR},
          #{email,jdbcType=VARCHAR},
          #{emailVerify,jdbcType=TINYINT},
          #{userMemCode,jdbcType=VARCHAR},
          #{createUser,jdbcType=VARCHAR},
          #{modifyDate,jdbcType=TIMESTAMP},
          #{modifyUser,jdbcType=VARCHAR})
  </insert>

    <select id="selectByConditions" resultMap="BaseResultMap">
        select
        `sysno`,
        `login_name`,
        `mobile`,
        `account_type`,
        `account_from_sys`,
        `account_status`,
        `invite_code`,
        `invite_account_sysno`,
        `nick_name`,
        `email`,
        `email_verify`,
        `user_memcode`,
        `create_date`,
        `create_user`,
        `modify_date`,
        `modify_user`
        from
        `account`
        <where>
            <if test="loginName != null and loginName != ''">
                and `login_name` like CONCAT('%', #{loginName,jdbcType=VARCHAR}, '%')
            </if>
            <if test="mobile != null and mobile != ''">
                and `mobile` like CONCAT('%', #{mobile,jdbcType=VARCHAR}, '%')
            </if>
            <if test="accountStatus != null ">
                and `account_status` = #{accountStatus}
            </if>
            <if test="accountType != null and accountType != ''">
                and `account_type` = #{accountType}
            </if>
            <if test="accountFromSys != null and accountFromSys != ''">
                and `account_from_sys` = #{accountFromSys}
            </if>
            <if test="email != null and email != ''">
                and `email` = #{email}
            </if>
            <if test="emailVerify != null and emailVerify != ''">
                and `email_verify` = #{emailVerify}
            </if>
            <if test="createUser != null and createUser != ''">
                and `create_user` = #{createUser}
            </if>
            <if test="modifyUser != null and modifyUser != ''">
                and `modify_user` = #{modifyUser}
            </if>
            <if test="startDate!=null and endDate!=null">
                <![CDATA[ AND `create_date` >= #{startDate,jdbcType=TIMESTAMP} AND `create_date` < #{endDate,jdbcType=TIMESTAMP} ]]>
            </if>
        </where>

    </select>

    <update id="updateByPrimaryKeySelective" parameterType="com.jwell.msa.customer.provider.account.entity.AccountDO">
        update
        `account`
        <set>
            <if test="loginName != null">
                `login_name` = #{loginName,jdbcType=VARCHAR},
            </if>
            <if test="mobile != null">
                `mobile` = #{mobile,jdbcType=VARCHAR},
            </if>
            <if test="accountType != null">
                `account_type` = #{accountType,jdbcType=TINYINT},
            </if>
            <if test="accountFromSys != null">
                `account_from_sys` = #{accountFromSys,jdbcType=VARCHAR},
            </if>
            <if test="accountStatus != null">
                `account_status` = #{accountStatus,jdbcType=TINYINT},
            </if>
            <if test="inviteCode != null">
                `invite_code` = #{inviteCode,jdbcType=VARCHAR},
            </if>
            <if test="inviteAccountSysno != null">
                `invite_account_sysno` = #{inviteAccountSysno,jdbcType=BIGINT},
            </if>
            <if test="nickName != null">
                `nick_name` = #{nickName,jdbcType=VARCHAR},
            </if>
            <if test="userMemCode != null">
                `user_memcode` = #{userMemCode,jdbcType=VARCHAR},
            </if>
            <if test="email != null and email != ''">
                `email` = #{email},
            </if>
            <if test="emailVerify != null and emailVerify != ''">
                `email_verify` = #{emailVerify},
            </if>
            <if test="modifyUser != null">
                `modify_user` = #{modifyUser,jdbcType=VARCHAR},
            </if>
        </set>
        where
        `sysno` = #{sysno,jdbcType=BIGINT}
    </update>

    <select id="selectAccountsByCompanySysno" resultMap="BaseResultMap">
      select
        m.*, mcr.administrator
      from
        account m,account_company_relation mcr
      where
        m.sysno = mcr.account_sysno
      and
        mcr.company_sysno= #{companySysno,jdbcType=BIGINT}
    </select>

    <select id="selectAccountByLoginName" resultMap="BaseResultMap">
      select
        `sysno`,
        `login_name`,
        `mobile`,
        `account_type`,
        `account_from_sys`,
        `account_status`,
        `invite_code`,
        `invite_account_sysno`,
        `nick_name`,
        `email`,
        `email_verify`,
        `user_memcode`,
        `create_date`,
        `create_user`,
        `modify_date`,
        `modify_user`
      from
        `account`
      where
        `login_name` = #{loginName,jdbcType=VARCHAR}
    </select>

    <update id="updateStatusByLoginName" parameterType="com.jwell.msa.customer.provider.account.entity.AccountDO">
        update
          `account`
        set
          `account_status` = #{accountStatus,jdbcType=TINYINT},
          `modify_user` = #{modifyUser,jdbcType=VARCHAR}
        where
          `login_name` = #{loginName,jdbcType=VARCHAR}
    </update>

    <select id="selectRegisterAccount" resultType="com.jwell.msa.customer.provider.account.entity.DayAccountDataDO">
        select count(1) as `count`, date_format(create_date,'%Y-%m-%d') as 'dayTime'
        from `account`
        where
        date_sub(curdate(), INTERVAL 6 DAY) &lt;= date(`create_date`)
        group by date_format(`create_date`,'%Y-%m-%d')
    </select>

    <select id="selectRegisterCount" resultType="int">
        select count(1)
        from `account` where
        <![CDATA[ `create_date` >= #{startDate,jdbcType=TIMESTAMP} AND `create_date` < #{endDate,jdbcType=TIMESTAMP}
        ]]>
    </select>

    <select id="selectAccountBySysnos" resultMap="BaseResultMap">
        select
        `sysno`,
        `login_name`,
        `mobile`,
        `account_type`,
        `account_from_sys`,
        `account_status`,
        `nick_name`
        from
        `account`
        where
        `sysno` in
        <foreach collection="list" index="index" item="item" open="(" separator="," close=")">
            #{item}
        </foreach>
    </select>

    <update id="updateMobilePrimaryKey"
            parameterType="com.jwell.msa.customer.provider.account.entity.AccountDO">
        update
        `account`
        set
        `mobile` = #{mobile,jdbcType=VARCHAR}
        where
        `sysno` = #{sysno,jdbcType=BIGINT}
    </update>

</mapper>