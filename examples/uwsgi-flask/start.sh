uwsgi --socket :4002 --protocol=http --workers=40 --module=app --callable=app