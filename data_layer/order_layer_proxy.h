#include <cstdint>
#include <string>

namespace order_layer_proxy {

using color_t = int8_t;
using seqnum_t = int64_t;

class IOrderClient {
   public:
    virtual void receive_order_response(seqnum_t token, seqnum_t gsn);
};

class OrderLayerProxy {
   public:
    OrderLayerProxy(std::string ip_addr, IOrderClient order_client);
    void send_order_request(color_t color, seqnum_t token);
};

}  // namespace order_layer_proxy