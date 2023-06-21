#include "tbb/concurrent_hash_map.h"
#include <cstdint>
#include <string>

using namespace tbb;

using LSNmap = concurrent_hash_map<uint64_t, std::string>;
using GSNmap = concurrent_hash_map<uint64_t, std::string>;

class CppLog {
private:
	LSNmap *lsnPtr;
    GSNmap *gsnPtr;
    uint64_t highest_gsn;
public:
    CppLog();
    ~CppLog();
    int Append(std::string record, uint64_t lsn);
    int Commit(uint64_t lsn, uint64_t gsn) ;
    const char *Read( uint64_t gsn, uint64_t *next_gsn);
    int Trim(uint64_t gsn);
};