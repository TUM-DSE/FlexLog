#pragma once

#include <cstdint>
#include <string>

#include "order_layer_proxy.h"

namespace data_layer {

using color_t = order_layer_proxy::color_t;
using seqnum_t = order_layer_proxy::seqnum_t;
using record_t = std::string;

// Interface for the log storage that stores one log for every color
// Has to be ***thread-safe***
class IStorage {
   public:
    // prepares the record by storing it at the end of the color's log,
    // findable with the token;
    // can return immediately, and actual storing can happen in the background,
    // important: the record cannnot be read after preparing
    virtual void prepare(color_t color, seqnum_t token, record_t record);

    // commits the record on the color's log that is stored finable at the token
    // with the global sequence number gsn,
    // waits for the record to be stored completely if still in prepation
    // the record can be read after the commit
    // throw exception if no record with lsn is stored or in preparation
    virtual void commit(seqnum_t token, seqnum_t gsn);

    // read's the record at the color's log at global sequence number gsn;
    // throw exception, if no such record
    virtual record_t read(color_t color, seqnum_t gsn);
};

struct AppendRequest {
    seqnum_t token;  // client-token, will be distinct across the whole shard
    color_t color;
    record_t record;
    bool responsible;  // if the replica is responsible to send the orderrequest
};

struct AppendResponse {
    seqnum_t gsn;
    // might want to add more later
};

struct ReadRequest {
    color_t color;
    seqnum_t gsn;
};

struct ReadResponse {
    record_t record;
    // might want to add more later
};

class Replica : order_layer_proxy::IOrderClient {
   public:
    Replica(IStorage storage, order_layer_proxy::OrderLayerProxy order_layer);

    AppendResponse append(AppendRequest app_req);
    ReadResponse read(ReadRequest read_req);

    void receive_order_response(seqnum_t token,
                                seqnum_t gsn) override;

   private:
    IStorage storage;
    order_layer_proxy::OrderLayerProxy order_layer;
};

}  // namespace data_layer