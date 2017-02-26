#include <ngx_config.h>
#include <ngx_core.h>
#include <ngx_http.h>

#define KEY_METHOD "Method"
#define KEY_REQUEST_URI "RequestURI"
#define KEY_PATH "Path"
#define KEY_ARGS "Args"
#define KEY_PROTO "Proto"
#define KEY_HOST "Host"
#define KEY_HEADERS "Headers"
#define KEY_HEADER_NAME "Name"
#define KEY_HEADER_VALUE "Value"

#define quoted_string(p, s, l) *p++ = '"'; \
                               p = (u_char *) ngx_escape_json(p, (u_char *) s, l); \
                               *p++ = '"'

#define quoted_string_len(s, l) sizeof("\"") - 1 \
                                + l + ngx_escape_json(NULL, (u_char *)s, l) \
                                + sizeof("\"") - 1
#ifndef ngx_copy_literal
#define ngx_copy_literal(p, s)  p = ngx_copy(p, s, sizeof(s) - 1)
#endif

static char *ngx_http_nose(ngx_conf_t *cf, ngx_command_t *cmd, void *conf);
static ngx_int_t ngx_http_nose_handler(ngx_http_request_t *r);
static uintptr_t ngx_http_nose_headers(ngx_http_request_t *r, u_char *dst);
static ngx_int_t ngx_http_nose_variable_host(ngx_http_request_t *r, ngx_http_variable_value_t *v);

static ngx_command_t ngx_http_nose_commands[] = {

    { ngx_string("http_nose_return"), /* directive */
      NGX_HTTP_LOC_CONF|NGX_CONF_NOARGS, /* location context and takes
                                            no arguments*/
      ngx_http_nose, /* configuration setup function */
      0, /* No offset. Only one context is supported. */
      0, /* No offset when storing the module configuration on struct. */
      NULL},

    ngx_null_command /* command termination */
};

/* The module context. */
static ngx_http_module_t ngx_http_nose_module_ctx = {
    NULL, /* preconfiguration */
    NULL, /* postconfiguration */

    NULL, /* create main configuration */
    NULL, /* init main configuration */

    NULL, /* create server configuration */
    NULL, /* merge server configuration */

    NULL, /* create location configuration */
    NULL /* merge location configuration */
};

/* Module definition. */
ngx_module_t ngx_http_nose_module = {
    NGX_MODULE_V1,
    &ngx_http_nose_module_ctx, /* module context */
    ngx_http_nose_commands, /* module directives */
    NGX_HTTP_MODULE, /* module type */
    NULL, /* init master */
    NULL, /* init module */
    NULL, /* init process */
    NULL, /* init thread */
    NULL, /* exit thread */
    NULL, /* exit process */
    NULL, /* exit master */
    NGX_MODULE_V1_PADDING
};

/**
 * Content handler.
 *
 * @param r
 *   Pointer to the request structure. See http_request.h.
 * @return
 *   The status of the response generation.
 */
static ngx_int_t ngx_http_nose_handler(ngx_http_request_t *r)
{
    size_t                      len;
    ngx_int_t                   rc;
    ngx_buf_t                   *b;
    ngx_chain_t                 *out;
    ngx_http_variable_value_t   *host;

    host = ngx_pcalloc(r->pool, sizeof(ngx_http_variable_value_t));
    if (host == NULL) {
        return NGX_ERROR;
    }
    ngx_http_nose_variable_host(r, host);

    /* calculate the buffer size */
    len = sizeof("{") - 1

        + quoted_string_len(KEY_METHOD, sizeof(KEY_METHOD) - 1)
        + sizeof(": ") - 1
        + quoted_string_len(r->method_name.data, r->method_name.len)
        + sizeof(", ") - 1

        + quoted_string_len(KEY_REQUEST_URI, sizeof(KEY_REQUEST_URI) - 1)
        + sizeof(": ") - 1
        + quoted_string_len(r->unparsed_uri.data, r->unparsed_uri.len)
        + sizeof(", ") - 1

        + quoted_string_len(KEY_PATH, sizeof(KEY_PATH) - 1)
        + sizeof(": ") - 1
        + quoted_string_len(r->uri.data, r->uri.len)
        + sizeof(", ") - 1

        + quoted_string_len(KEY_ARGS, sizeof(KEY_ARGS) - 1)
        + sizeof(": ") - 1
        + quoted_string_len(r->args.data, r->args.len)
        + sizeof(", ") - 1

        + quoted_string_len(KEY_PROTO, sizeof(KEY_PROTO) - 1)
        + sizeof(": ") - 1
        + quoted_string_len(r->http_protocol.data, r->http_protocol.len)
        + sizeof(", ") - 1
        ;
    if (!host->not_found) {
        len +=
            + quoted_string_len(KEY_HOST, sizeof(KEY_HOST) - 1)
            + sizeof(": ") - 1
            + quoted_string_len(host->data, host->len)
            + sizeof(", ") - 1
            ;
    }
    len +=
        + quoted_string_len(KEY_HEADERS, sizeof(KEY_HEADERS) - 1)
        + sizeof(": ") - 1
        + sizeof("[") - 1
        + ngx_http_nose_headers(r, NULL)
        + sizeof("]") - 1

        + sizeof("}") - 1
        ;

    /* create the buffer */

    out = ngx_alloc_chain_link(r->pool);
    if (out == NULL) {
        return NGX_ERROR;
    }

    b = ngx_create_temp_buf(r->pool, len);
    if (b == NULL) {
        return NGX_ERROR;
    }
    /* there will be no more buffers in the request */
    b->last_buf = 1;

    out->buf = b;
    out->next = NULL;

    /* copy data over to the buffer */
    *b->last++ = '{';

    quoted_string(b->last, KEY_METHOD, sizeof(KEY_METHOD) - 1);
    ngx_copy_literal(b->last, ": ");
    quoted_string(b->last, r->method_name.data, r->method_name.len);
    ngx_copy_literal(b->last, ", ");

    quoted_string(b->last, KEY_REQUEST_URI, sizeof(KEY_REQUEST_URI) - 1);
    ngx_copy_literal(b->last, ": ");
    quoted_string(b->last, r->unparsed_uri.data, r->unparsed_uri.len);
    ngx_copy_literal(b->last, ", ");

    quoted_string(b->last, KEY_PATH, sizeof(KEY_PATH) - 1);
    ngx_copy_literal(b->last, ": ");
    quoted_string(b->last, r->uri.data, r->uri.len);
    ngx_copy_literal(b->last, ", ");

    quoted_string(b->last, KEY_ARGS, sizeof(KEY_ARGS) - 1);
    ngx_copy_literal(b->last, ": ");
    quoted_string(b->last, r->args.data, r->args.len);
    ngx_copy_literal(b->last, ", ");

    quoted_string(b->last, KEY_PROTO, sizeof(KEY_PROTO) - 1);
    ngx_copy_literal(b->last, ": ");
    quoted_string(b->last, r->http_protocol.data, r->http_protocol.len);
    ngx_copy_literal(b->last, ", ");

    if (!host->not_found) {
        quoted_string(b->last, KEY_HOST, sizeof(KEY_HOST) - 1);
        ngx_copy_literal(b->last, ": ");
        quoted_string(b->last, host->data, host->len);
        ngx_copy_literal(b->last, ", ");
    }

    quoted_string(b->last, KEY_HEADERS, sizeof(KEY_HEADERS) - 1);
    ngx_copy_literal(b->last, ": ");
    *b->last++ = '[';
    b->last = (u_char *) ngx_http_nose_headers(r, b->last);
    *b->last++ = ']';

    *b->last++ = '}';

    if (b->last != b->end) {
        ngx_log_error(NGX_LOG_ERR, r->connection->log, 0,
                "http_nose: buffer error");

        return NGX_ERROR;
    }

    // /* Sending the headers for the reply. */

    r->headers_out.status = NGX_HTTP_OK;
    r->headers_out.content_length_n = len;
    r->headers_out.content_type.len = sizeof("application/json") - 1;
    r->headers_out.content_type.data = (u_char *) "application/json";;

    rc = ngx_http_send_header(r);
    if (rc == NGX_ERROR || rc >= NGX_HTTP_SPECIAL_RESPONSE) {
        return rc;
    }

    return ngx_http_output_filter(r, out);
} /* ngx_http_nose_handler */

/**
 * Configuration setup function that installs the content handler.
 *
 * @param cf
 *   Module configuration structure pointer.
 * @param cmd
 *   Module directives structure pointer.
 * @param conf
 *   Module configuration structure pointer.
 * @return string
 *   Status of the configuration setup.
 */
static char *ngx_http_nose(ngx_conf_t *cf, ngx_command_t *cmd, void *conf)
{
    ngx_http_core_loc_conf_t *clcf; /* pointer to core location configuration */

    clcf = ngx_http_conf_get_module_loc_conf(cf, ngx_http_core_module);
    clcf->handler = ngx_http_nose_handler;

    return NGX_CONF_OK;
} /* ngx_http_nose */

/**
 * Add request headers.
*/
static uintptr_t
ngx_http_nose_headers(ngx_http_request_t *r, u_char *dst)
{
    uintptr_t        len;
    ngx_uint_t        i;
    ngx_list_part_t  *part;
    ngx_table_elt_t  *header;

    part = &r->headers_in.headers.part;
    header = part->elts;

    if (dst == NULL)
        len = 0;

    for (i = 0; /* void */ ; i++) {

        if (i >= part->nelts) {
            if (part->next == NULL) {
                break;
            }

            part = part->next;
            header = part->elts;
            i = 0;
        }

        if (header[i].hash == 0) {
            continue;
        }

        if (dst == NULL) {
            len += sizeof("{") - 1

                + quoted_string_len(KEY_HEADER_NAME, sizeof(KEY_HEADER_NAME) - 1)
                + sizeof(": ") - 1
                + quoted_string_len(header[i].key.data, header[i].key.len)
                + sizeof(", ") - 1

                + quoted_string_len(KEY_HEADER_VALUE, sizeof(KEY_HEADER_VALUE) - 1)
                + sizeof(": ") - 1
                + quoted_string_len(header[i].value.data, header[i].value.len)

                + sizeof("}") - 1
                ;

            if (i + 1 < part->nelts)
                len += sizeof(", ") - 1;

        } else {
                *dst++ = '{';
                quoted_string(dst, KEY_HEADER_NAME, sizeof(KEY_HEADER_NAME) - 1);
                ngx_copy_literal(dst, ": ");
                quoted_string(dst, header[i].key.data, header[i].key.len);
                ngx_copy_literal(dst, ", ");

                quoted_string(dst, KEY_HEADER_VALUE, sizeof(KEY_HEADER_VALUE) - 1);
                ngx_copy_literal(dst, ": ");
                quoted_string(dst, header[i].value.data, header[i].value.len);
                *dst++ = '}';

                if (i + 1 < part->nelts)
                    ngx_copy_literal(dst, ", ");
        }
    }
    if (dst == NULL)
        return (uintptr_t) len;
    return (uintptr_t) dst;
} /* ngx_http_nose_headers */

static ngx_int_t
ngx_http_nose_variable_host(ngx_http_request_t *r, ngx_http_variable_value_t *v)
{
    ngx_http_core_srv_conf_t  *cscf;

    if (r->headers_in.server.len) {
        v->len = r->headers_in.server.len;
        v->data = r->headers_in.server.data;

    } else {
        cscf = ngx_http_get_module_srv_conf(r, ngx_http_core_module);

        v->len = cscf->server_name.len;
        v->data = cscf->server_name.data;
    }

    v->valid = 1;
    v->no_cacheable = 0;
    v->not_found = 0;

    return NGX_OK;
}