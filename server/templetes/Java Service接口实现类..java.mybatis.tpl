package {{ project.NameSpace }}.service;

import com.dk.foundation.engine.baseentity.PageResult;
import com.dk.foundation.engine.exception.BusinessException;
import {{ project.NameSpace }}.dao.{{ fn.ToCamelString(table.Name) }}Mapper;
import {{ project.NameSpace }}.entity.{{ fn.ToCamelString(table.Name) }};
import {{ project.NameSpace }}.req.{{ fn.ToCamelString(table.Name) }}Filter;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import com.github.pagehelper.Page;
import com.github.pagehelper.PageHelper;

import java.util.List;

@Service
public class {{ fn.ToCamelString(table.Name) }}ServiceImpl implements {{ fn.ToCamelString(table.Name) }}Service {

    @Autowired
    private {{ fn.ToCamelString(table.Name) }}Mapper {{ fn.FirstToLower(fn.ToCamelString(table.Name)) }}Mapper;

    @Override
    public {{ fn.SqlTypeConvertLanguageType(fn.GetPK(indexs),project.DataBase,project.Language) }} insert({{ fn.ToCamelString(table.Name) }} entity) throws BusinessException {
        if (entity == null) {
            throw new BusinessException("参数不能为空");
        }
        {{ fn.FirstToLower(fn.ToCamelString(table.Name)) }}Mapper.insert(entity);
        return entity.get{{ fn.ToCamelString(fn.GetPK(indexs).Name) }}();
    }

    @Override
    public void update({{ fn.ToCamelString(table.Name) }} entity) throws BusinessException {
        if (entity == null) {
            throw new BusinessException("参数不能为空");
        }

        {{ fn.FirstToLower(fn.ToCamelString(table.Name)) }}Mapper.update(entity);
    }

    @Override
    public void delete({{ fn.SqlTypeConvertLanguageType(fn.GetPK(indexs),project.DataBase,project.Language) }} {{ fn.FirstToLower(fn.ToCamelString(fn.GetPK(indexs).Name)) }}) {
        {{ fn.FirstToLower(fn.ToCamelString(table.Name)) }}Mapper.delete({{ fn.FirstToLower(fn.ToCamelString(fn.GetPK(indexs).Name)) }});
    }

    @Override
    public {{ fn.ToCamelString(table.Name) }} load({{ fn.SqlTypeConvertLanguageType(fn.GetPK(indexs),project.DataBase,project.Language) }} {{ fn.FirstToLower(fn.ToCamelString(fn.GetPK(indexs).Name)) }}) {
        return {{ fn.FirstToLower(fn.ToCamelString(table.Name)) }}Mapper.load({{ fn.FirstToLower(fn.ToCamelString(fn.GetPK(indexs).Name)) }});
    }

    @Override
    public PageResult<{{ fn.ToCamelString(table.Name) }}> queryByPage({{ fn.ToCamelString(table.Name) }}Filter filter) {
        Page page = PageHelper.startPage(filter.getPageIndex(), filter.getPageSize(), filter.getOrderBy());
        List<{{ fn.ToCamelString(table.Name) }}> result= {{ fn.FirstToLower(fn.ToCamelString(table.Name)) }}Mapper.query(filter);
        return new PageResult<>(result,page.getTotal(),filter.getPageIndex(),filter.getPageSize());
    }

}