(self.webpackChunkdoc_ops=self.webpackChunkdoc_ops||[]).push([[836],{7781:function(e,t,n){"use strict";n.r(t),n.d(t,{frontMatter:function(){return s},contentTitle:function(){return c},metadata:function(){return u},toc:function(){return l},default:function(){return h}});var r=n(2122),o=n(9756),i=(n(7294),n(3905)),a=["components"],s={},c="FAQ",u={unversionedId:"faq",id:"faq",isDocsHomePage:!1,title:"FAQ",description:"What is GoShimmer?",source:"@site/docs/faq.md",sourceDirName:".",slug:"/faq",permalink:"/docs/faq",editUrl:"https://github.com/iotaledger/goshimmer/tree/develop/documentation/docs/faq.md",version:"current",frontMatter:{},sidebar:"docs",previous:{title:"Welcome",permalink:"/docs/welcome"},next:{title:"Setting up a GoShimmer node",permalink:"/docs/tutorials/setup"}},l=[{value:"What is GoShimmer?",id:"what-is-goshimmer",children:[]},{value:"What Kind of Confirmation Time Can I Expect?",id:"what-kind-of-confirmation-time-can-i-expect",children:[]},{value:"Where Can I See the State of the GoShimmer testnet?",id:"where-can-i-see-the-state-of-the-goshimmer-testnet",children:[]},{value:"How Many Transactions Per Second(TPS) can GoShimmer Sustain?",id:"how-many-transactions-per-secondtps-can-goshimmer-sustain",children:[]},{value:"How is Spamming Prevented?",id:"how-is-spamming-prevented",children:[]},{value:"What Happens if I Issue a Double Spend?",id:"what-happens-if-i-issue-a-double-spend",children:[]},{value:"Who&#39;s the Target Audience for Operating a GoShimmer Node?",id:"whos-the-target-audience-for-operating-a-goshimmer-node",children:[]}],d={toc:l};function h(e){var t=e.components,n=(0,o.Z)(e,a);return(0,i.kt)("wrapper",(0,r.Z)({},d,n,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"faq"},"FAQ"),(0,i.kt)("h2",{id:"what-is-goshimmer"},"What is GoShimmer?"),(0,i.kt)("p",null,"GoShimmer is a research and engineering project from the IOTA Foundation seeking to evaluate Coordicide concepts by implementing them in a node software."),(0,i.kt)("h2",{id:"what-kind-of-confirmation-time-can-i-expect"},"What Kind of Confirmation Time Can I Expect?"),(0,i.kt)("p",null,"Since non-conflicting transactions aren't even voted on, they materialize after 2x the average network delay parameter we set. This means that a transaction usually confirms within a time boundary of ~10 seconds."),(0,i.kt)("h2",{id:"where-can-i-see-the-state-of-the-goshimmer-testnet"},"Where Can I See the State of the GoShimmer testnet?"),(0,i.kt)("p",null,"You can access the global analysis dashboard in the ",(0,i.kt)("a",{parentName:"p",href:"http://ressims.iota.cafe:28080/autopeering"},"Pollen Analyzer")," which showcases the network graph and active ongoing votes on conflicts."),(0,i.kt)("h2",{id:"how-many-transactions-per-secondtps-can-goshimmer-sustain"},"How Many Transactions Per Second(TPS) can GoShimmer Sustain?"),(0,i.kt)("p",null,"The transactions per second metric is irrelevant for the current development state of GoShimmer. We are evaluating components from Coordicide, and aren't currently interested in squeezing out every little ounce of performance. Since the primary goal is to evaluate Coordicide components, we value simplicity over optimization . Even if we would put out a TPS number, it would not reflect an actual metric in a finished production ready node software. "),(0,i.kt)("h2",{id:"how-is-spamming-prevented"},"How is Spamming Prevented?"),(0,i.kt)("p",null,"The Coordicide lays out concepts for spam prevention through the means of rate control and such. However, in the current version, GoShimmer relies on Proof of Work (PoW) to prevent over saturation of the network. Doing the PoW for a message will usually take a couple of seconds on commodity hardware."),(0,i.kt)("h2",{id:"what-happens-if-i-issue-a-double-spend"},"What Happens if I Issue a Double Spend?"),(0,i.kt)("p",null,"If issue simultaneous transactions spending the same funds, there is high certainty that your transaction will be rejected by the network. This rejection will block your funds indefinitely, though this may change in the future.  "),(0,i.kt)("p",null,"If you issue a transaction, await the average network delay, and then issue the double spend, then the first issued transaction should usually become confirmed, and the 2nd one rejected.  "),(0,i.kt)("h2",{id:"whos-the-target-audience-for-operating-a-goshimmer-node"},"Who's the Target Audience for Operating a GoShimmer Node?"),(0,i.kt)("p",null,"Our primary focus is testing out Coordicide components. We are mainly interested in individuals who have a strong IT background, rather than giving people of any knowledge-level the easiest way to operate a node. We welcome people interested in trying out the bleeding edge of IOTA development and providing meaningful feedback or problem reporting in form of ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/iotaledger/goshimmer/issues/new/choose"},"issues"),"."))}h.isMDXComponent=!0},3905:function(e,t,n){"use strict";n.d(t,{Zo:function(){return l},kt:function(){return m}});var r=n(7294);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function a(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function s(e,t){if(null==e)return{};var n,r,o=function(e,t){if(null==e)return{};var n,r,o={},i=Object.keys(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var c=r.createContext({}),u=function(e){var t=r.useContext(c),n=t;return e&&(n="function"==typeof e?e(t):a(a({},t),e)),n},l=function(e){var t=u(e.components);return r.createElement(c.Provider,{value:t},e.children)},d={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},h=r.forwardRef((function(e,t){var n=e.components,o=e.mdxType,i=e.originalType,c=e.parentName,l=s(e,["components","mdxType","originalType","parentName"]),h=u(n),m=o,p=h["".concat(c,".").concat(m)]||h[m]||d[m]||i;return n?r.createElement(p,a(a({ref:t},l),{},{components:n})):r.createElement(p,a({ref:t},l))}));function m(e,t){var n=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var i=n.length,a=new Array(i);a[0]=h;var s={};for(var c in t)hasOwnProperty.call(t,c)&&(s[c]=t[c]);s.originalType=e,s.mdxType="string"==typeof e?e:o,a[1]=s;for(var u=2;u<i;u++)a[u]=n[u];return r.createElement.apply(null,a)}return r.createElement.apply(null,n)}h.displayName="MDXCreateElement"}}]);