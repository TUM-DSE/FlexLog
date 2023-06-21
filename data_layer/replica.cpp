#include "data_layer.h"

namespace data_layer {

Replica::Replica(IStorage storage,
                 order_layer_proxy::OrderLayerProxy order_layer)
    : storage{storage}, order_layer{order_layer} {}

AppendResponse Replica::append(AppendRequest app_req) {
    // TODO
}

ReadResponse Replica::read(ReadRequest read_req) {
    // TODO
}

void Replica::receive_order_response(seqnum_t token,
                                     seqnum_t gsn) {
    // TODO
}

}  // namespace data_layer