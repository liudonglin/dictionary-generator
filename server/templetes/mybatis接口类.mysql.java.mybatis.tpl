package {{ project.NameSpace }}.dao;

import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;
import java.util.List;

/**
 * {{ table.Description }}
 * @author {{ user }}
 * @date {{ fn.GetNowDate() }}
 * */
@Mapper
public interface {{ fn.ToCamelString(table.Name) }}Mapper {

    /**
     * 新增{{ table.Description }}
     * */
    Integer insert({{ fn.ToCamelString(table.Name) }} entity);

    /**
     * 修改{{ table.Description }}
     * */
    Integer update({{ fn.ToCamelString(table.Name) }} entity);

    /**
     * 根据主键加载{{ table.Description }}
     * */
    {{ fn.ToCamelString(table.Name) }} load(Long id);

    /**
     * 查询{{ table.Description }}
     * */
    List<{{ fn.ToCamelString(table.Name) }}> query({{ fn.ToCamelString(table.Name) }}Filter filter);

    /**
     * 根据主键删除{{ table.Description }}
     * */
    void delete(Long id);
}
