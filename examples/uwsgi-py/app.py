import json

def application(environ, start_response):
    start_response('200 OK', [('Content-Type', 'application/json')])
    headers = []
    for k, v in environ.items():
        if k.startswith('HTTP_'):
            headers.append({
                'Name': k[5:].lower().replace('_', '-'),
                'Value': v
            })
    data = {'Headers': headers}
    if 'REQUEST_METHOD' in environ:
        data['Method'] = environ['REQUEST_METHOD']
    if 'REQUEST_URI' in environ:
        data['RequestURI'] = environ['REQUEST_URI']
    if 'PATH_INFO' in environ:
        data['Path'] = environ['PATH_INFO']
    if 'QUERY_STRING' in environ:
        data['Args'] = environ['QUERY_STRING']
    if 'SERVER_PROTOCOL' in environ:
        data['Proto'] = environ['SERVER_PROTOCOL']
    if 'SERVER_NAME' in environ:
        data['Host'] = environ['SERVER_NAME']
    if environ.get('CONTENT_LENGTH', None):
        data['Body'] = environ['wsgi.input'].read().decode('utf-8')

    return [json.dumps(data, indent=2)]
