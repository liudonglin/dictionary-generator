package {{ project.NameSpace }}.service;

import {{ project.NameSpace }}.util.PageResult;
import {{ project.NameSpace }}.util.BusinessException;
import {{ project.NameSpace }}.req.{{ fn.ToCamelString(table.Name) }}Filter;

import java.util.List;
import java.util.Map;

/**
 * @author: {{ user.Login }}
 * @date: {{ fn.GetNowDate() }}
 * @description: {{ table.Description }}
 */

public interface {{ fn.ToCamelString(table.Name) }}Service {

    /**
     * 添加{{ table.Description }}
     * @param entity
     * @return
     */
    {{ fn.SqlTypeConvertLanguageType(fn.GetPK(indexs),project.DataBase,project.Language) }} insert({{ fn.ToCamelString(table.Name) }} entity) throws BusinessException;

    /**
     * 更新{{ table.Description }}
     * @param entity
     * @return
     */
    void update({{ fn.ToCamelString(table.Name) }} entity) throws BusinessException;

    /**
     * 删除{{ table.Description }}
     * @param entity
     * @return
     */
    void delete({{ fn.SqlTypeConvertLanguageType(fn.GetPK(indexs),project.DataBase,project.Language) }} {{ fn.FirstToLower(fn.ToCamelString(fn.GetPK(indexs).Name)) }});

    /**
     * 加载{{ table.Description }}
     * @param {{ fn.FirstToLower(fn.ToCamelString(fn.GetPK(indexs).Name)) }}
     * @return
     */
    {{ fn.ToCamelString(table.Name) }} load({{ fn.SqlTypeConvertLanguageType(fn.GetPK(indexs),project.DataBase,project.Language) }} {{ fn.FirstToLower(fn.ToCamelString(fn.GetPK(indexs).Name)) }});

    /**
     * 分页查询{{ table.Description }}
     * @param {{ fn.ToCamelString(table.Name) }}Filter
     * @return
     */
    PageResult<{{ fn.ToCamelString(table.Name) }}> selectByPage({{ fn.ToCamelString(table.Name) }}Filter filter);

}