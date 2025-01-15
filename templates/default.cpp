#include <bits/stdc++.h>
using namespace std;

template<typename A, typename B> ostream& operator<<(ostream &os, const pair<A, B> &p) { return os << '(' << p.first << ", " << p.second << ')'; }
template<typename T_container, typename T = typename enable_if<!is_same<T_container, string>::value, typename T_container::value_type>::type> ostream& operator<<(ostream &os, const T_container &v) { os << '{'; string sep; for (const T &x : v) os << sep << x, sep = ", "; return os << '}'; }
void dbg_out() { cerr << endl; }
template<typename Head, typename... Tail> void dbg_out(Head H, Tail... T) { cerr << ' ' << H; dbg_out(T...); }


#define ar array
#define ll long long
#define ld long double
#define sza(x) ((int)x.size())
#define all(a) (a).begin(), (a).end()

#ifdef YASH_DEBUG
#define debug(x) cerr << __LINE__ << ": " << #x << " "; __print(x); cerr << endl;
#else
#define debug(x)
#endif
 
template<typename T> void __print(vector<T> vec) {
    cerr << "[ ";
    for(auto it: vec) cerr << it << " ";
    cerr << "]";
}
template<typename T> void __print(set<T> vec) {
    cerr << "[ ";
    for(auto it: vec) cerr << it << " ";
    cerr << "]";
}
template<typename T, typename T2> void __print(map<T, T2> m){
    cerr << "[ ";
    for(auto it: m){
        cerr << "{ " << it.first << "," <<it.second<<" }, ";
    }
    cerr << "]";
}

template<typename T, typename T2> void __print(unordered_map<T, T2> m){
    cerr << "[ ";
    for(auto it: m){
        cerr << "{ " << it.first << "," <<it.second<<" }, ";
    }
    cerr << "]";
}
void __print(ll a) {cerr << a;}
void __print(ld a) {cerr << a;}
void __print(char a) {cerr << a;}
void __print(string a) {cerr << a;}
void __print(int a) {cerr << a;}
void __print(float a){cerr << a;}


const int MAX_N = 1e5 + 5;
const ll MOD = 1e9 + 7;
const ll INF = 1e9;
const ld EPS = 1e-9;



void solve() {

}

int main() {
#ifdef YASH_DEBUG
    freopen("Debug.txt", "w", stderr);
    cerr << __FILE__ << endl;
#endif
    ios_base::sync_with_stdio(0);
    cin.tie(0); cout.tie(0);
    int tc = 1;
    // cin >> tc;
    for (int t = 1; t <= tc; t++) {
        // cout << "Case #" << t << ": ";
        solve();
    }
}
