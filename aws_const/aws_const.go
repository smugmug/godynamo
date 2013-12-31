// Copyright (c) 2013,2014 SmugMug, Inc. All rights reserved.
// 
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//     * Redistributions of source code must retain the above copyright
//       notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
//       copyright notice, this list of conditions and the following
//       disclaimer in the documentation and/or other materials provided
//       with the distribution.
// 
// THIS SOFTWARE IS PROVIDED BY SMUGMUG, INC. ``AS IS'' AND ANY
// EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
// PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL SMUGMUG, INC. BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE
// GOODS OR SERVICES;LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER
// IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR
// OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
// ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// A colelction of consts reused in various packages.
package aws_const

const (
	DYNAMODB                 = "DynamoDB"
	ISO8601FMT               = "2006-01-02T15:04:05Z"
	ISO8601FMT_CONDENSED     = "20060102T150405Z"
	ISODATEFMT               = "20060102"
	PORT                     = "80"
	METHOD                   = "POST"
	CTYPE                    = "application/x-amz-json-1.0"
	AMZ_TARGET_HDR           = "X-Amz-Target"
	CONTENT_MD5_HDR          = "Content-MD5"
	CONTENT_TYPE_HDR         = "Content-Type"
	DATE_HDR                 = "Date"
	CURRENT_API_VERSION      = "DynamoDB_20120810"
	ENDPOINT_PREFIX          = CURRENT_API_VERSION + "."
	X_AMZ_DATE_HDR           = "X-Amz-Date"
	X_AMZ_SECURITY_TOKEN_HDR = "X-Amz-Security-Token"
	X_AMZN_AUTHORIZATION_HDR = "X-Amzn-Authorization"
	RETRIES                  = 7
	EXCEEDED_MSG             = "ProvisionedThroughputExceededException"
	UNRECOGNIZED_CLIENT_MSG  = "UnrecognizedClientException"
	THROTTLING_MSG           = "ThrottlingException"
)
