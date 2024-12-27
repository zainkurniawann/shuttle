-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "route_jawa" (
	"route_id" BIGINT PRIMARY KEY,
	"route_uuid" UUID NOT NULL,
	"driver_uuid" UUID NOT NULL,
	"student_uuid" UUID NOT NULL,
	"route_name" VARCHAR(100) NOT NULL,
	"route_description" TEXT NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NOT NULL,
	"created_by" VARCHAR(255) NOT NULL,
	"updated_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"updated_by" VARCHAR(255) NULL DEFAULT NULL,
	"deleted_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"deleted_by" VARCHAR(255) NULL DEFAULT NULL,
	"school_uuid" UUID NOT NULL,
	FOREIGN KEY ("driver_uuid") REFERENCES "driver_details" ("user_uuid") ON UPDATE NO ACTION ON DELETE NO ACTION,
	FOREIGN KEY ("driver_uuid") REFERENCES "users" ("user_uuid") ON UPDATE NO ACTION ON DELETE NO ACTION,
	FOREIGN KEY ("driver_uuid") REFERENCES "driver_details" ("user_uuid") ON UPDATE NO ACTION ON DELETE SET NULL,
	FOREIGN KEY ("driver_uuid") REFERENCES "users" ("user_uuid") ON UPDATE NO ACTION ON DELETE SET NULL,
	FOREIGN KEY ("school_uuid") REFERENCES "schools" ("school_uuid") ON UPDATE NO ACTION ON DELETE NO ACTION,
	FOREIGN KEY ("student_uuid") REFERENCES "students" ("student_uuid") ON UPDATE NO ACTION ON DELETE NO ACTION
);

INSERT INTO "route_jawa" ("route_id", "route_uuid", "driver_uuid", "student_uuid", "route_name", "route_description", "created_at", "created_by", "updated_at", "updated_by", "deleted_at", "deleted_by", "school_uuid") VALUES
	(1735212264937943865, '406885f3-2865-4095-bb3d-9b4081ac4df1', '5ad039d9-33e2-4331-a757-372c325fdc56', '1aa23748-8981-4be4-bc02-611e2d07b543', 'apalah kocak', 'homes to school in the morning.', '2024-12-26 18:24:24.937306+07', 'admin2ngaglik', NULL, NULL, NULL, NULL, '5b659fa0-1d68-459f-afe6-b59e8f0e4f97'),
	(1735212307896049121, '9b9c0c6e-1392-4707-bcbf-35f065a3226b', '5ad039d9-33e2-4331-a757-372c325fdc56', '33f7a0f2-be3f-4ecc-9aeb-868d7fe90aa8', 'apalah kocak2', 'homes to school in the morning.', '2024-12-26 18:25:07.896319+07', 'admin2ngaglik', NULL, NULL, NULL, NULL, '5b659fa0-1d68-459f-afe6-b59e8f0e4f97');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS route_jawa;
-- +goose StatementEnd