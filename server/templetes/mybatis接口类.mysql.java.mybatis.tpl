package {{ project.NameSpace }}.dao;

import {{ project.NameSpace }}.req.{{ fn.ToCamelString(table.Name) }}Filter;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;
import java.util.List;

/**
 * @author: {{ user.Login }}
 * @date: {{ fn.GetNowDate() }}
 * @description: {{ table.Description }}
 */
@Mapper
public interface {{ fn.ToCamelString(table.Name) }}Mapper {

    /**
     * 新增{{ table.Description }}
     */
    Integer insert({{ fn.ToCamelString(table.Name) }} entity);

    /**
     * 修改{{ table.Description }}
     */
    Integer update({{ fn.ToCamelString(table.Name) }} entity);

    /**
     * 根据主键加载{{ table.Description }}
     */
    {{ fn.ToCamelString(table.Name) }} load({{ fn.SqlTypeConvertLanguageType(fn.GetPK(indexs),project.DataBase,project.Language) }} {{ fn.FirstToLower(fn.ToCamelString(fn.GetPK(indexs).Name)) }});

    /**
     * 查询{{ table.Description }}
     */
    List<{{ fn.ToCamelString(table.Name) }}> query({{ fn.ToCamelString(table.Name) }}Filter filter);

    /**
     * 根据主键删除{{ table.Description }}
     */
    void delete({{ fn.SqlTypeConvertLanguageType(fn.GetPK(indexs),project.DataBase,project.Language) }} {{ fn.FirstToLower(fn.ToCamelString(fn.GetPK(indexs).Name)) }});
}
