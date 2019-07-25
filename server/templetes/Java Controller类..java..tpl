package {{ project.NameSpace }}.controller;

import com.dk.foundation.engine.baseentity.PageResult;
import com.dk.foundation.engine.baseentity.StandResponse;
import com.dk.foundation.engine.exception.BusinessException;
import {{ project.NameSpace }}.entity.FlowType;
import {{ project.NameSpace }}.req.FlowTypeFilter;
import {{ project.NameSpace }}.service.FlowTypeService;
import {{ project.NameSpace }}.util.BaseController;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiImplicitParam;
import io.swagger.annotations.ApiOperation;
import org.springframework.stereotype.Controller;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

import javax.annotation.Resource;

@Controller
@RequestMapping("/v1/{{ fn.ToLower(fn.ToCamelString(table.Name)) }}")
@Api(value = "{{ table.Title }}", produces = "application/json;charset=UTF-8")
@Validated
public class {{ fn.ToCamelString(table.Name) }}Controller extends BaseController {

    @Resource
    private {{ fn.ToCamelString(table.Name) }}Service {{ fn.FirstToLower(fn.ToCamelString(table.Name)) }}Service;

    @ApiOperation(value = "添加{{ table.Title }}", notes = "")
    @ApiImplicitParam(name = "Authorization", value = "获取token后，直接传入", paramType = "header")
    @RequestMapping(value = "/create", method = RequestMethod.POST)
    public @ResponseBody StandResponse create{{ fn.ToCamelString(table.Name) }}(@RequestBody {{ fn.ToCamelString(table.Name) }} entity) throws BusinessException {
        if (entity==null) {
            throw new BusinessException("参数不能为空");
        }
        {{ fn.FirstToLower(fn.ToCamelString(table.Name)) }}Service.insert(entity);
        return success(entity.getId());
    }

    @ApiOperation(value = "修改{{ table.Title }}", notes = "")
    @ApiImplicitParam(name = "Authorization", value = "获取token后，直接传入", paramType = "header")
    @RequestMapping(value = "/update", method = RequestMethod.POST)
    public @ResponseBody StandResponse update{{ fn.ToCamelString(table.Name) }}(@RequestBody {{ fn.ToCamelString(table.Name) }} entity) throws BusinessException {
        if (entity==null) {
            throw new BusinessException("参数不能为空");
        }
        {{ fn.FirstToLower(fn.ToCamelString(table.Name)) }}Service.update(entity);
        return success();
    }

    @ApiOperation(value = "删除{{ table.Title }}", notes = "")
    @ApiImplicitParam(name = "Authorization", value = "获取token后，直接传入", paramType = "header")
    @RequestMapping(value = "/delete", method = RequestMethod.POST)
    public @ResponseBody StandResponse delete{{ fn.ToCamelString(table.Name) }}(@RequestParam {{ fn.SqlTypeConvertLanguageType(fn.GetPK(indexs),project.DataBase,project.Language) }} id) throws BusinessException {
        if (id==null) {
            throw new BusinessException("参数不能为空");
        }
        flowTypeService.delete(id);
        return success();
    }

    @ApiOperation(value = "根据主键获取{{ table.Title }}", notes = "")
    @ApiImplicitParam(name = "Authorization", value = "获取token后，直接传入", paramType = "header")
    @RequestMapping(value = "/load", method = RequestMethod.POST)
    public @ResponseBody StandResponse<{{ fn.ToCamelString(table.Name) }}> load{{ fn.ToCamelString(table.Name) }}(@RequestParam {{ fn.SqlTypeConvertLanguageType(fn.GetPK(indexs),project.DataBase,project.Language) }} id) throws BusinessException {
        if (id==null) {
            throw new BusinessException("参数不能为空");
        }
        return success({{ fn.FirstToLower(fn.ToCamelString(table.Name)) }}Service.load(id));
    }

    @ApiOperation(value = "分页查询{{ table.Title }}", notes = "")
    @ApiImplicitParam(name = "Authorization", value = "获取token后，直接传入", paramType = "header")
    @RequestMapping(value = "/query", method = RequestMethod.POST)
    public @ResponseBody StandResponse<PageResult<{{ fn.ToCamelString(table.Name) }}>> query{{ fn.ToCamelString(table.Name) }}(@RequestBody {{ fn.ToCamelString(table.Name) }}Filter req) throws BusinessException {
        if (req==null) {
            throw new BusinessException("参数不能为空");
        }
        return success({{ fn.FirstToLower(fn.ToCamelString(table.Name)) }}Service.queryByPage(req));
    }

}
