package driver

/**
 * Drivers must be able to create the table below.
 *
 * TABLE_NAME: follow `schemaManagementTableName` constant
 * COLUMNS:
 *   version: not null, unique, can assign 64 bits unsigned integer
 */

const schemaManagementTableName = "schema_version"
