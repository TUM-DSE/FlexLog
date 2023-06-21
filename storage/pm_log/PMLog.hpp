#include <libpmemobj++/persistent_ptr.hpp>
#include <libpmemobj++/container/concurrent_hash_map.hpp>
#include <cstdint>
#include <string>

using namespace pmem::obj;

struct root;

class PersistentString {
	private:
		persistent_ptr<char[]> array;
	public:
		const char* data() const;	
		PersistentString(const std::string&);
};

using PString = struct PersistentString;
using LSNmap = concurrent_hash_map<p<uint64_t>, persistent_ptr<PString>>;
using GSNmap = concurrent_hash_map<p<uint64_t>, persistent_ptr<PString>>;

class cppPMLog {
	private:
		persistent_ptr<LSNmap> lsnPptr;
		persistent_ptr<GSNmap> gsnPptr;	
		p<uint64_t> highest_gsn;
		p<uint64_t> lowest_gsn;
		void cacheRecords(cppPMLog *log, uint64_t gsn);
	public:
		p<pool<root>> pop;
		cppPMLog(pool<root>);
		~cppPMLog();
		void restartMaps();
		void shutdown();
		int Append(const char *record, uint64_t lsn);
		int Commit(uint64_t lsn, uint64_t gsn) ;
		const char *Read(uint64_t gsn, uint64_t *next_gsn);
		int Trim(uint64_t gsn);
};

void setup(std::string &s1, std::string &s2);
void* cppStartUp();
void* cppStartUp(int);
void cppFinalize(persistent_ptr<cppPMLog> cppLog);
