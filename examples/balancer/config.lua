instance = {
  addrs = {
   { ip = "localhost"; port = 8000; };
  };
    http = {
        maxlen = 65536; maxreq = 65536;
        headers = {
            create_func = {
                ["X-Func"] = "realip"; ["Host"] = "host";
            };

            create_func_weak = {
                ["X-Func-Weak"] = "realip";
            };

            proxy = {
                host = "localhost"; port = 9000;
            }; -- proxy
        }; -- headers
    }; --http

}; --instance
