DBNAME:=annotation
DBINST:=video_db
MIGRATEPATH:=datastore/migrations

mock:
	@mockgen -source=./datastore/repository.go -destination=./rest_services/mock.go -package=rest_services
migrate_sql:
	migrate create -ext sql -dir ${MIGRATEPATH} -seq  annotation_schema

createdb:
	docker exec -it ${DBINST} createdb --username=root --owner=root ${DBNAME}
	#docker exec -it ${DBINST} psgl -U root ${DBNAME}

dropdb:
	docker exec -it ${DBINST} dropdb ${DBNAME}

migratedown:
	migrate -path ${MIGRATEPATH} -database "$(DB_URL)" -verbose down
migrateup:
	migrate -path ${MIGRATEPATH} -database "$(DB_URL)" -verbose up

.PHONY: mock migrate_sql createdb dropdb migrateup migratedown