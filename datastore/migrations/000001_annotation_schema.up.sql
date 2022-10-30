-- noinspection SqlNoDataSourceInspectionForFile

CREATE TABLE "video_schema" (
                                "id" bigserial PRIMARY KEY,
                                "video" varchar UNIQUE,
                                "duration" float,
                                "frame_rate" float,
                                "schema" jsonb,
                                "created_at" timestamptz DEFAULT (now()),
                                "updated_at" timestamptz
);

CREATE TABLE "annotation" (
                           "id" bigserial PRIMARY KEY,
                           "video_schema_id" bigserial,
                           "start" float,
                           "end_time" float,
                           "metadata" jsonb,
                           "created_at" timestamptz DEFAULT (now()),
                           "updated_at" timestamptz
);


ALTER TABLE "annotation" ADD FOREIGN KEY ("video_schema_id") REFERENCES "video_schema" ("id");

