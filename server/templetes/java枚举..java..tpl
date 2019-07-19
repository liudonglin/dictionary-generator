package {{ project.NameSpace }}.entity;

{% for enum in enums %}
/**
* {{ enum.Title }}
*/
public enum {{ enum.Name }} {
    {% for kvd in enum.Kvds %}
    /**
     * {{ kvd.Des }}
     */
    {{ fn.ToUpper(kvd.Key) }}({{ kvd.Value }}, "{{ kvd.Des }}"){% if fn.IsLastKvd(kvd,enum.Kvds) %};{% else %},{% endif %}{% endfor %}

    /**
     * 代码
     */
    private final Integer code;
    /**
     * 信息
     */
    private final String message;

    private {{ enum.Name }}(Integer code, String message) {
        this.code = code;
        this.message = message;
    }

    public static {{ enum.Name }} getByCode(Integer code) {
        if (code == null) {
            return null;
        }
        for ({{ enum.Name }} value : values()) {
            if (code.equals(value.getCode())) {
                return value;
            }
        }
        return null;
    }

    public Integer getCode() {
        return code;
    }

    public String getMessage() {
        return message;
    }

}
{% endfor %}