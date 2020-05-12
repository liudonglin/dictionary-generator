package {{ project.NameSpace }}.service;

import com.oursoft.scf.base.PageResult;
import com.oursoft.scf.base.BusinessException;
import {{ project.NameSpace }}.entity.{{ fn.ToCamelString(table.Name) }};
import {{ project.NameSpace }}.req.{{ fn.ToCamelString(table.Name) }}Filter;

/**
 * @author: {{ user.Login }}
 * @date: {{ fn.GetNowDate() }}
 * @description: {{ table.Title }}
 */

public interface {{ fn.ToCamelString(table.Name) }}Service {

    /**
     * 添加{{ table.Title }}
     * @param entity
     * @return
     */
    {{ fn.SqlTypeConvertLanguageType(fn.GetPK(indexs),project.DataBase,project.Language) }} insert({{ fn.ToCamelString(table.Name) }} entity) throws BusinessException;

    /**
     * 更新{{ table.Title }}
     * @param entity
     * @return
     */
    void update({{ fn.ToCamelString(table.Name) }} entity) throws BusinessException;

    /**
     * 删除{{ table.Title }}
     * @param {{ fn.FirstToLower(fn.ToCamelString(fn.GetPK(indexs).Name)) }}
     * @return
     */
    void delete({{ fn.SqlTypeConvertLanguageType(fn.GetPK(indexs),project.DataBase,project.Language) }} {{ fn.FirstToLower(fn.ToCamelString(fn.GetPK(indexs).Name)) }});

    /**
     * 加载{{ table.Title }}
     * @param {{ fn.FirstToLower(fn.ToCamelString(fn.GetPK(indexs).Name)) }}
     * @return
     */
    {{ fn.ToCamelString(table.Name) }} load({{ fn.SqlTypeConvertLanguageType(fn.GetPK(indexs),project.DataBase,project.Language) }} {{ fn.FirstToLower(fn.ToCamelString(fn.GetPK(indexs).Name)) }});

    /**
     * 分页查询{{ table.Title }}
     * @param filter
     * @return
     */
    PageResult<{{ fn.ToCamelString(table.Name) }}> queryByPage({{ fn.ToCamelString(table.Name) }}Filter filter);

}