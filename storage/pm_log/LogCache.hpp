#include "tbb/concurrent_hash_map.h"
#include <cstdint>
#include <string>
#include <memory>

using GSNmapCache = tbb::concurrent_hash_map<uint64_t, std::string>;

class LogCache {
	private:
		std::shared_ptr<GSNmapCache> logCachePtr;
		uint64_t highest_gsn;
		uint64_t lowest_gsn;
		uint64_t next_gsn(uint64_t gsn);
	public:
		LogCache();
		~LogCache();
		int Append(const std::string record, uint64_t gsn);
		const char *Read(uint64_t gsn, uint64_t *next_gsn);
		int Erase(uint64_t gsn);
};
