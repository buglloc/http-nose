from flask import Flask, request, Response
import json

app = Flask(__name__)

@app.before_request
def catch_all():
    args = request.path.split('?', 1)
    args = args[1] if len(args) == 2 else ''
    headers = []
    for k, v in request.headers.items():
        headers.append({
            'Name': k,
            'Value': v
        })

    data = {
        'Method': request.method,
        'RequestURI': request.full_path,
        'Path': request.path,
        'Args': args,
        'Headers': headers,
        'Body': request.data.decode('utf-8'),
        'Host': request.headers['Host']
    }
    return Response(json.dumps(data, indent=2),  mimetype='application/json')


if __name__ == '__main__':
    app.run(port=4000)
