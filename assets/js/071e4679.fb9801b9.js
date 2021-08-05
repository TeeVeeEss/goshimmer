(self.webpackChunkdoc_ops=self.webpackChunkdoc_ops||[]).push([[5774],{1590:function(e,t,n){"use strict";n.r(t),n.d(t,{frontMatter:function(){return r},contentTitle:function(){return l},metadata:function(){return c},toc:function(){return d},default:function(){return h}});var a=n(2122),i=n(9756),o=(n(7294),n(3905)),s=["components"],r={},l="Congestion Control",c={unversionedId:"protocol_specification/congestion_control",id:"protocol_specification/congestion_control",isDocsHomePage:!1,title:"Congestion Control",description:"Every network has to deal with its intrinsic limited resources in terms of bandwidth and node capabilities (CPU and storage). In this document, we present a congestion control algorithm to regulate the influx of messages in the network with the goal of maximizing throughput (messages/bytes per second) and minimizing delays. Furthermore, the following requirements must be satisfied:",source:"@site/docs/protocol_specification/congestion_control.md",sourceDirName:"protocol_specification",slug:"/protocol_specification/congestion_control",permalink:"/docs/protocol_specification/congestion_control",editUrl:"https://github.com/iotaledger/goshimmer/tree/develop/documentation/docs/protocol_specification/congestion_control.md",version:"current",frontMatter:{},sidebar:"docs",previous:{title:"Mana Implementation",permalink:"/docs/protocol_specification/mana"},next:{title:"Consensus Mechanism",permalink:"/docs/protocol_specification/consensus_mechanism"}},d=[{value:"Detailed design",id:"detailed-design",children:[{value:"Prerequirements",id:"prerequirements",children:[]},{value:"Outbox management",id:"outbox-management",children:[]},{value:"Scheduler",id:"scheduler",children:[]},{value:"Rate setting",id:"rate-setting",children:[]},{value:"Message blocking and blacklisting",id:"message-blocking-and-blacklisting",children:[]}]}],u={toc:d};function h(e){var t=e.components,r=(0,i.Z)(e,s);return(0,o.kt)("wrapper",(0,a.Z)({},u,r,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("h1",{id:"congestion-control"},"Congestion Control"),(0,o.kt)("p",null,"Every network has to deal with its intrinsic limited resources in terms of bandwidth and node capabilities (CPU and storage). In this document, we present a congestion control algorithm to regulate the influx of messages in the network with the goal of maximizing throughput (messages/bytes per second) and minimizing delays. Furthermore, the following requirements must be satisfied:"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("em",{parentName:"li"},"Consistency"),". If a message is written by one honest node, it shall be written by all honest nodes within some delay bound."),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("em",{parentName:"li"},"Fairness"),". Nodes can obtain a share of the available throughput depending on their access Mana. Throughput is shared in such a way that an attempt to increase the allocation of any node necessarily results in the decrease in the allocation of some other node with an equal or smaller allocation (max-min fairness)."),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("em",{parentName:"li"},"Security"),". Malicious nodes shall be unable to interfere with either of the above requirements.")),(0,o.kt)("p",null,(0,o.kt)("img",{alt:"Congestion Control",src:n(948).Z})),(0,o.kt)("p",null,"Further information can be found in the paper ",(0,o.kt)("a",{parentName:"p",href:"https://arxiv.org/abs/2005.07778"},"Access Control for Distributed Ledgers in the Internet of Things: A Networking Approach"),"."),(0,o.kt)("h2",{id:"detailed-design"},"Detailed design"),(0,o.kt)("p",null,"Our algorithm has three core components: "),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"A scheduling algorithm which ensures fair access for all nodes according to their access Mana."),(0,o.kt)("li",{parentName:"ul"},"A TCP-inspired algorithm for decentralized rate setting to efficiently utilize the available bandwidth while preventing large delays."),(0,o.kt)("li",{parentName:"ul"},"A blacklisting policy to ban malicious nodes.")),(0,o.kt)("h3",{id:"prerequirements"},"Prerequirements"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("p",{parentName:"li"},(0,o.kt)("em",{parentName:"p"},"Node identity"),". We require node accountability where each message is associated with the node ID of its issuing node.")),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("p",{parentName:"li"},(0,o.kt)("em",{parentName:"p"},"Access mana"),". The congestion control module has knowledge of the access Mana of the nodes in the network in order to fairly share the available throughput. Without access Mana the network would be subject to Sybil attacks, which would incentivise even honest actors to artificially increase its own number of nodes.")),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("p",{parentName:"li"},(0,o.kt)("em",{parentName:"p"},"Timestamp"),". Before scheduling a new message, the scheduler verifies whether the message timestamp is valid or not.")),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("p",{parentName:"li"},(0,o.kt)("em",{parentName:"p"},"Message weight"),". Weight of a message is used to priority messages over the others and it is calculated depending on the type of message and of the message length."))),(0,o.kt)("h3",{id:"outbox-management"},"Outbox management"),(0,o.kt)("p",null,"Once the message has successfully passed the message parser checks and is solid, it is enqueued into the outbox for scheduling. The outbox is logically split into several queues, each one corresponding to a different node issuing messages. In this section, we describe the operations of message enqueuing (and dequeuing) into (from) the outbox."),(0,o.kt)("p",null,"The enqueuing mechanism includes the following components:"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("em",{parentName:"li"},"Classification"),". The mechanism identifies the queue where the message belongs to according to the node ID of the message issuer."),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("em",{parentName:"li"},"Message enqueuing"),". The message is actually enqueued, queue is sorted by message timestamps in increasing order and counters are updated (e.g., counters for the total number of bytes in the queue)."),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("em",{parentName:"li"},"Message drop"),". In some circumstances, due to network congestion or to ongoing attacks, some messages shall be dropped to guarantee bounded delays and isolate attacker's messages. Specifically, a node shall drop messages in two situations:",(0,o.kt)("ul",{parentName:"li"},(0,o.kt)("li",{parentName:"ul"},"since buffers are of a limited size, if the total number of bytes in all queues exceeds a certain threshold, new incoming messages are dropped;"),(0,o.kt)("li",{parentName:"ul"},"to guarantee the security of the network, if a certain queue exceeds a given threshold, new incoming packets from that specific node ID will be dropped.")))),(0,o.kt)("p",null,"The dequeue mechanism includes the following components:"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("em",{parentName:"li"},"Queue selection"),". A queue is selected according to round robin scheduling algorithm. In particular, we use a modified version of the deficit round robin (DRR) algorithm."),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("em",{parentName:"li"},"Message dequeuing"),". The first message of the queue is dequeued, and list of active nodes is updated."),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("em",{parentName:"li"},"Scheduler management"),". Scheduler counters and pointers are updated.")),(0,o.kt)("h3",{id:"scheduler"},"Scheduler"),(0,o.kt)("p",null,"The most critical task is the scheduling algorithm which must guarantee that, for an honest node ",(0,o.kt)("inlineCode",{parentName:"p"},"node"),", the following requirements will be met:"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("inlineCode",{parentName:"li"},"node"),"'s messages will not accumulate indefinitely at any node (i.e., starvation is avoided), so the ",(0,o.kt)("em",{parentName:"li"},"consistency")," requirement will be ensured."),(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("inlineCode",{parentName:"li"},"node"),"'s fair share (according to its access Mana) of the network resources are allocated to it, guaranteeing the ",(0,o.kt)("em",{parentName:"li"},"fairness")," requirement."),(0,o.kt)("li",{parentName:"ul"},"Malicious nodes sending above their allowed rate will not interrupt ",(0,o.kt)("inlineCode",{parentName:"li"},"node"),"'s throughput, fulfilling the ",(0,o.kt)("em",{parentName:"li"},"security")," requirement.")),(0,o.kt)("p",null,"Although nodes in our setting are capable of more complex and customised behaviour than a typical router in a packet-switched network, our scheduler must still be lightweight and scalable due to the potentially large number of nodes requiring differentiated treatment. It is estimated that over 10,000 nodes operate on the Bitcoin network, and we expect that an even greater number of nodes are likely to be present in the IoT setting. For this reason, we adopt a scheduler based on ",(0,o.kt)("a",{parentName:"p",href:"https://ieeexplore.ieee.org/document/502236"},"Deficit Round Robin")," (DRR) (the Linux implementation of the ",(0,o.kt)("a",{parentName:"p",href:"https://tools.ietf.org/html/rfc8290"},"FQ-CoDel packet scheduler"),", which is based on DRR, supports anywhere up to 65535 separate queues)."),(0,o.kt)("p",null,"The DRR scans all non-empty queues in sequence. When a non-empty queue is selected, its priority counter (called ",(0,o.kt)("em",{parentName:"p"},"deficit"),") is incremented by a certain value (called ",(0,o.kt)("em",{parentName:"p"},"quantum"),"). Then, the value of the deficit counter is a maximal amount of bytes that can be sent at this turn: if the deficit counter is greater than the weight of the message at the head of the queue, this message can be scheduled and the value of the counter is decremented by this weight. In our implementation, the quantum is proportional to node's access Mana and we add a cap on the maximum deficit that a node can achieve to keep the network latency low. It is also important to mention that the weight of the message can be assigned in such a way that specific messages can be prioritized (low weight) or penalized (large weight); by default, in our mechanism the weight is proportional to the message size measured in bytes. The weight of a message is set by the function ",(0,o.kt)("inlineCode",{parentName:"p"},"WorkCalculator()"),"."),(0,o.kt)("p",null,"Here a fundamental remark: ",(0,o.kt)("em",{parentName:"p"},"the network manager sets up a desired maximum (fixed) rate")," ",(0,o.kt)("inlineCode",{parentName:"p"},"SCHEDULING_RATE")," ",(0,o.kt)("em",{parentName:"p"},"at which messages will be scheduled"),", computed in weight (see above) per second. This implies that every message is scheduled after a delay which is equal to the weight (size as default) of the latest scheduled message times the parameter ",(0,o.kt)("inlineCode",{parentName:"p"},"SCHEDULING_RATE"),". This rate mostly depends on the degree of decentralization desired: e.g., a larger rate leads to higher throughput but would leave behind slower devices which will fall out of sync."),(0,o.kt)("h3",{id:"rate-setting"},"Rate setting"),(0,o.kt)("p",null,"If all nodes always had messages to issue, i.e., if nodes were continuously willing to issue new messages, the problem of rate setting would be very straightforward: nodes could simply operate at a fixed, assured rate, sharing the total throughput according to the percentage of access Mana owned. The scheduling algorithm would ensure that this rate is enforceable, and that increasing delays or dropped messages are only experienced by misbehaving node. However, it is unrealistic that all nodes will always have messages to issue, and we would like nodes to better utilise network resources, without causing excessive congestion and violating any requirement."),(0,o.kt)("p",null,"We propose a rate setting algorithm inspired by TCP \u2014 each node employs ",(0,o.kt)("a",{parentName:"p",href:"https://https://epubs.siam.org/doi/book/10.1137/1.9781611974225"},"additive increase, multiplicative decrease")," (AIMD) rules to update their issuance rate in response to congestion events. In the case of distributed ledgers, all message traffic passes through all nodes, contrary to the case of traffic typically found in packet switched networks and other traditional network architectures. Under these conditions, local congestion at a node is all that is required to indicate congestion elsewhere in the network. This observation is crucial, as it presents an opportunity for a congestion control algorithm based entirely on local traffic."),(0,o.kt)("p",null,"Our rate setting algorithm outlines the AIMD rules employed by each node to set their issuance rate. Rate updates for a node ",(0,o.kt)("inlineCode",{parentName:"p"},"node")," take place each time a new message is scheduled if the ",(0,o.kt)("inlineCode",{parentName:"p"},"node")," has a non-empty set of its own messages not yet scheduled. Node ",(0,o.kt)("inlineCode",{parentName:"p"},"node")," sets its own local additive-increase variable ",(0,o.kt)("inlineCode",{parentName:"p"},"localIncrease(node)")," based on its access Mana and on a global increase rate parameter ",(0,o.kt)("inlineCode",{parentName:"p"},"RATE_SETTING_INCREASE"),". An appropriate choice of ",(0,o.kt)("inlineCode",{parentName:"p"},"RATE_SETTING_INCREASE")," ensures a conservative global increase rate which does not cause problems even when many nodes increase their rate simultaneously. Nodes wait ",(0,o.kt)("inlineCode",{parentName:"p"},"RATE_SETTING_PAUSE")," seconds after a global multiplicative decrease parameter ",(0,o.kt)("inlineCode",{parentName:"p"},"RATE_SETTING_DECREASE"),", during which there are no further updates made, to allow the reduced rate to take effect and prevent multiple successive decreases. At each update, ",(0,o.kt)("inlineCode",{parentName:"p"},"node")," checks how many of its own messages are in its outbox queue, and responds with a multiplicative decrease if this number is above a threshold, ",(0,o.kt)("inlineCode",{parentName:"p"},"backoff(node)"),", which is proportional to ",(0,o.kt)("inlineCode",{parentName:"p"},"node"),"'s access Mana. If the number of ",(0,o.kt)("inlineCode",{parentName:"p"},"node"),"'s messages in the outbox is below the threshold, ",(0,o.kt)("inlineCode",{parentName:"p"},"node"),"'s issuance rate is incremented by its local increase variable ",(0,o.kt)("inlineCode",{parentName:"p"},"localIncrease(node)"),"."),(0,o.kt)("h3",{id:"message-blocking-and-blacklisting"},"Message blocking and blacklisting"),(0,o.kt)("p",null,"If an incoming message made the outbox total buffer size to exceed its maximum capacity ",(0,o.kt)("inlineCode",{parentName:"p"},"MAX_BUFFER"),", the same message would be dropped. In our analysis, we set buffers to be large enough to accommodate traffic from all honest nodes."),(0,o.kt)("p",null,"Furthermore, to mitigate spamming actions from malicious nodes, we add an additional constraint: if ",(0,o.kt)("inlineCode",{parentName:"p"},"node"),"'s access Mana-scaled queue length (i.e., queue length divided by node's access Mana) exceeds a given threshold ",(0,o.kt)("inlineCode",{parentName:"p"},"MAX_QUEUE"),", any new incoming packet from ",(0,o.kt)("inlineCode",{parentName:"p"},"node")," will be dropped, hence the node is blacklisted. The attacker is blacklisted for a certain time ",(0,o.kt)("inlineCode",{parentName:"p"},"BLACKLIST_TIME")," during which no messages issued by ",(0,o.kt)("inlineCode",{parentName:"p"},"node")," can be added to the outbox. Please note that it is still possible to receive message from the attacker through solidification requests, which is important in order to guarantee the consistency requirement. Finally, when a node is blacklisted, the blacklister does not increase its own rate for a time ",(0,o.kt)("inlineCode",{parentName:"p"},"RATE_SETTING_QUARANTINE"),", to avoid errors in the perception of the current congestion level."))}h.isMDXComponent=!0},3905:function(e,t,n){"use strict";n.d(t,{Zo:function(){return d},kt:function(){return m}});var a=n(7294);function i(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function s(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){i(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function r(e,t){if(null==e)return{};var n,a,i=function(e,t){if(null==e)return{};var n,a,i={},o=Object.keys(e);for(a=0;a<o.length;a++)n=o[a],t.indexOf(n)>=0||(i[n]=e[n]);return i}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(a=0;a<o.length;a++)n=o[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(i[n]=e[n])}return i}var l=a.createContext({}),c=function(e){var t=a.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):s(s({},t),e)),n},d=function(e){var t=c(e.components);return a.createElement(l.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return a.createElement(a.Fragment,{},t)}},h=a.forwardRef((function(e,t){var n=e.components,i=e.mdxType,o=e.originalType,l=e.parentName,d=r(e,["components","mdxType","originalType","parentName"]),h=c(n),m=i,p=h["".concat(l,".").concat(m)]||h[m]||u[m]||o;return n?a.createElement(p,s(s({ref:t},d),{},{components:n})):a.createElement(p,s({ref:t},d))}));function m(e,t){var n=arguments,i=t&&t.mdxType;if("string"==typeof e||i){var o=n.length,s=new Array(o);s[0]=h;var r={};for(var l in t)hasOwnProperty.call(t,l)&&(r[l]=t[l]);r.originalType=e,r.mdxType="string"==typeof e?e:i,s[1]=r;for(var c=2;c<o;c++)s[c]=n[c];return a.createElement.apply(null,s)}return a.createElement.apply(null,n)}h.displayName="MDXCreateElement"},948:function(e,t,n){"use strict";t.Z=n.p+"assets/images/congestion_control_algorithm_infographic-2e9a9d99c4bf3de4c3c980e3bff5d969.png"}}]);