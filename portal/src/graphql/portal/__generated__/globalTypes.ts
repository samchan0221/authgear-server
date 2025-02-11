/* tslint:disable */
/* eslint-disable */
// @generated
// This file was automatically generated and should not be edited.

//==============================================================
// START Enums and Input Objects
//==============================================================

/**
 * Update to resource file.
 */
export interface AppResourceUpdate {
  data?: string | null;
  path: string;
}

export interface OauthClientSecretInput {
  alias: string;
  clientSecret: string;
}

export interface SecretConfigInput {
  oauthClientSecrets?: OauthClientSecretInput[] | null;
}

//==============================================================
// END Enums and Input Objects
//==============================================================
