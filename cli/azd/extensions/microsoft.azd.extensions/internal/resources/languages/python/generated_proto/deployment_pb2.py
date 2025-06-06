# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: deployment.proto
# Protobuf Python Version: 5.29.0
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    29,
    0,
    '',
    'deployment.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


import models_pb2 as models__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x10\x64\x65ployment.proto\x12\x06\x61zdext\x1a\x0cmodels.proto\"?\n\x15GetDeploymentResponse\x12&\n\ndeployment\x18\x01 \x01(\x0b\x32\x12.azdext.Deployment\"J\n\x1cGetDeploymentContextResponse\x12*\n\x0c\x41zureContext\x18\x01 \x01(\x0b\x32\x14.azdext.AzureContext\"\xaa\x02\n\nDeployment\x12\n\n\x02id\x18\x01 \x01(\t\x12\x10\n\x08location\x18\x02 \x01(\t\x12\x14\n\x0c\x64\x65ploymentId\x18\x03 \x01(\t\x12\x0c\n\x04name\x18\x04 \x01(\t\x12\x0c\n\x04type\x18\x05 \x01(\t\x12*\n\x04tags\x18\x06 \x03(\x0b\x32\x1c.azdext.Deployment.TagsEntry\x12\x30\n\x07outputs\x18\x07 \x03(\x0b\x32\x1f.azdext.Deployment.OutputsEntry\x12\x11\n\tresources\x18\x08 \x03(\t\x1a+\n\tTagsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\x1a.\n\x0cOutputsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\r\n\x05value\x18\x02 \x01(\t:\x02\x38\x01\x32\xad\x01\n\x11\x44\x65ploymentService\x12\x44\n\rGetDeployment\x12\x14.azdext.EmptyRequest\x1a\x1d.azdext.GetDeploymentResponse\x12R\n\x14GetDeploymentContext\x12\x14.azdext.EmptyRequest\x1a$.azdext.GetDeploymentContextResponseBFZ4github.com/azure/azure-dev/cli/azd/pkg/azdext;azdext\xaa\x02\rMicrosoft.Azdb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'deployment_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z4github.com/azure/azure-dev/cli/azd/pkg/azdext;azdext\252\002\rMicrosoft.Azd'
  _globals['_DEPLOYMENT_TAGSENTRY']._loaded_options = None
  _globals['_DEPLOYMENT_TAGSENTRY']._serialized_options = b'8\001'
  _globals['_DEPLOYMENT_OUTPUTSENTRY']._loaded_options = None
  _globals['_DEPLOYMENT_OUTPUTSENTRY']._serialized_options = b'8\001'
  _globals['_GETDEPLOYMENTRESPONSE']._serialized_start=42
  _globals['_GETDEPLOYMENTRESPONSE']._serialized_end=105
  _globals['_GETDEPLOYMENTCONTEXTRESPONSE']._serialized_start=107
  _globals['_GETDEPLOYMENTCONTEXTRESPONSE']._serialized_end=181
  _globals['_DEPLOYMENT']._serialized_start=184
  _globals['_DEPLOYMENT']._serialized_end=482
  _globals['_DEPLOYMENT_TAGSENTRY']._serialized_start=391
  _globals['_DEPLOYMENT_TAGSENTRY']._serialized_end=434
  _globals['_DEPLOYMENT_OUTPUTSENTRY']._serialized_start=436
  _globals['_DEPLOYMENT_OUTPUTSENTRY']._serialized_end=482
  _globals['_DEPLOYMENTSERVICE']._serialized_start=485
  _globals['_DEPLOYMENTSERVICE']._serialized_end=658
# @@protoc_insertion_point(module_scope)
