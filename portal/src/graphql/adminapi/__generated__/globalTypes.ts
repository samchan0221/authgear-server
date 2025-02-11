/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

//==============================================================
// START Enums and Input Objects
//==============================================================

export enum AuditLogActivityType {
  AUTHENTICATION_IDENTITY_ANONYMOUS_FAILED = "AUTHENTICATION_IDENTITY_ANONYMOUS_FAILED",
  AUTHENTICATION_IDENTITY_BIOMETRIC_FAILED = "AUTHENTICATION_IDENTITY_BIOMETRIC_FAILED",
  AUTHENTICATION_IDENTITY_LOGIN_ID_FAILED = "AUTHENTICATION_IDENTITY_LOGIN_ID_FAILED",
  AUTHENTICATION_PRIMARY_OOB_OTP_EMAIL_FAILED = "AUTHENTICATION_PRIMARY_OOB_OTP_EMAIL_FAILED",
  AUTHENTICATION_PRIMARY_OOB_OTP_SMS_FAILED = "AUTHENTICATION_PRIMARY_OOB_OTP_SMS_FAILED",
  AUTHENTICATION_PRIMARY_PASSWORD_FAILED = "AUTHENTICATION_PRIMARY_PASSWORD_FAILED",
  AUTHENTICATION_SECONDARY_OOB_OTP_EMAIL_FAILED = "AUTHENTICATION_SECONDARY_OOB_OTP_EMAIL_FAILED",
  AUTHENTICATION_SECONDARY_OOB_OTP_SMS_FAILED = "AUTHENTICATION_SECONDARY_OOB_OTP_SMS_FAILED",
  AUTHENTICATION_SECONDARY_PASSWORD_FAILED = "AUTHENTICATION_SECONDARY_PASSWORD_FAILED",
  AUTHENTICATION_SECONDARY_RECOVERY_CODE_FAILED = "AUTHENTICATION_SECONDARY_RECOVERY_CODE_FAILED",
  AUTHENTICATION_SECONDARY_TOTP_FAILED = "AUTHENTICATION_SECONDARY_TOTP_FAILED",
  IDENTITY_EMAIL_ADDED = "IDENTITY_EMAIL_ADDED",
  IDENTITY_EMAIL_REMOVED = "IDENTITY_EMAIL_REMOVED",
  IDENTITY_EMAIL_UPDATED = "IDENTITY_EMAIL_UPDATED",
  IDENTITY_OAUTH_CONNECTED = "IDENTITY_OAUTH_CONNECTED",
  IDENTITY_OAUTH_DISCONNECTED = "IDENTITY_OAUTH_DISCONNECTED",
  IDENTITY_PHONE_ADDED = "IDENTITY_PHONE_ADDED",
  IDENTITY_PHONE_REMOVED = "IDENTITY_PHONE_REMOVED",
  IDENTITY_PHONE_UPDATED = "IDENTITY_PHONE_UPDATED",
  IDENTITY_USERNAME_ADDED = "IDENTITY_USERNAME_ADDED",
  IDENTITY_USERNAME_REMOVED = "IDENTITY_USERNAME_REMOVED",
  IDENTITY_USERNAME_UPDATED = "IDENTITY_USERNAME_UPDATED",
  USER_ANONYMOUS_PROMOTED = "USER_ANONYMOUS_PROMOTED",
  USER_AUTHENTICATED = "USER_AUTHENTICATED",
  USER_CREATED = "USER_CREATED",
  USER_SIGNED_OUT = "USER_SIGNED_OUT",
}

export enum AuthenticatorKind {
  PRIMARY = "PRIMARY",
  SECONDARY = "SECONDARY",
}

export enum AuthenticatorType {
  OOB_OTP_EMAIL = "OOB_OTP_EMAIL",
  OOB_OTP_SMS = "OOB_OTP_SMS",
  PASSWORD = "PASSWORD",
  TOTP = "TOTP",
}

export enum IdentityType {
  ANONYMOUS = "ANONYMOUS",
  BIOMETRIC = "BIOMETRIC",
  LOGIN_ID = "LOGIN_ID",
  OAUTH = "OAUTH",
}

export enum SessionType {
  IDP = "IDP",
  OFFLINE_GRANT = "OFFLINE_GRANT",
}

export enum SortDirection {
  ASC = "ASC",
  DESC = "DESC",
}

export enum UserSortBy {
  CREATED_AT = "CREATED_AT",
  LAST_LOGIN_AT = "LAST_LOGIN_AT",
}

/**
 * Definition of an identity. This is a union object, exactly one of the available fields must be present.
 */
export interface IdentityDefinition {
  loginID?: IdentityDefinitionLoginID | null;
}

export interface IdentityDefinitionLoginID {
  key: string;
  value: string;
}

//==============================================================
// END Enums and Input Objects
//==============================================================
