<?xml version="1.0" encoding="UTF-8" ?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd" >
<mapper namespace="{{ project.NameSpace }}.dao.{{ fn.ToCamelString(table.Name) }}Mapper">
    
    <select id="load" resultType="{{ project.NameSpace }}.entity.{{ fn.ToCamelString(table.Name) }}">
        SELECT {% for column in columns %}
        {{column.Name}}{% if !fn.IsLastColumn(column,columns) %},{% endif %}{% endfor %}
        FROM {{ table.Name }}
        WHERE {{ fn.GetPK(indexs).Name }} = #{ {{ fn.FirstToLower(fn.ToCamelString(fn.GetPK(indexs).Name)) }} }
    </select>

    <insert id="insert" parameterType="{{ project.NameSpace }}.entity.{{ fn.ToCamelString(table.Name) }}" {% if fn.GetPK(indexs).Name!="" %} keyProperty="{{ fn.FirstToLower(fn.ToCamelString(fn.GetPK(indexs).Name)) }}" useGeneratedKeys="true" {% endif %}>
        INSERT INTO {{ table.Name }} ( {% for column in columns %} {% if !column.AI %}
        {{column.Name}}{% if !fn.IsLastColumn(column,columns) %},{% endif %}{% endif %}{% endfor %}
        )
        values ( {% for column in columns %} {% if !column.AI %}
        #{ {{fn.FirstToLower(fn.ToCamelString(column.Name))}} }{% if !fn.IsLastColumn(column,columns) %},{% endif %}{% endif %}{% endfor %}
        )
    </insert>

    <update id="update" parameterType="{{ project.NameSpace }}.entity.{{ fn.ToCamelString(table.Name) }}">
        UPDATE {{ table.Name }}
        <set> {% for column in columns %} {% if !column.PK %}
            <if test="{{fn.FirstToLower(fn.ToCamelString(column.Name))}} != null">
                {{column.Name}} = #{ {{fn.FirstToLower(fn.ToCamelString(column.Name))}} }{% if !fn.IsLastColumn(column,columns) %},{% endif %}
            </if> {% endif %}{% endfor %}
        </set>
        WHERE {{ fn.GetPK(indexs).Name }} = #{ {{ fn.FirstToLower(fn.ToCamelString(fn.GetPK(indexs).Name)) }} }
    </update>

    <delete id="delete" >
        DELETE FROM {{ table.Name }} WHERE {{ fn.GetPK(indexs).Name }} = #{ {{ fn.FirstToLower(fn.ToCamelString(fn.GetPK(indexs).Name)) }} }
    </delete>

    <select id="query" resultType="{{ project.NameSpace }}.entity.{{ fn.ToCamelString(table.Name) }}" parameterType="{{ project.NameSpace }}.req.{{ fn.ToCamelString(table.Name) }}Filter">
        SELECT {% for column in columns %}
        {{column.Name}}{% if !fn.IsLastColumn(column,columns) %},{% endif %}{% endfor %}
        FROM {{ table.Name }}
        <where> {% for column in columns %}
            <if test="{{fn.FirstToLower(fn.ToCamelString(column.Name))}} != null">
                and {{column.Name}} {% if fn.SqlTypeConvertLanguageType(column,project.DataBase,project.Language)=="String" %} like CONCAT ('%', #{ {{fn.FirstToLower(fn.ToCamelString(column.Name))}} },'%') {% else %}= #{ {{fn.FirstToLower(fn.ToCamelString(column.Name))}} }{% endif %}
            </if>{% endfor %}
        </where>
    </select>

</mapper>