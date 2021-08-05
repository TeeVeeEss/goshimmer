(self.webpackChunkdoc_ops=self.webpackChunkdoc_ops||[]).push([[5138],{5835:function(e,t,n){"use strict";n.r(t),n.d(t,{frontMatter:function(){return s},contentTitle:function(){return l},metadata:function(){return c},toc:function(){return p},default:function(){return d}});var r=n(2122),o=n(9756),i=(n(7294),n(3905)),a=["components"],s={},l="Integration tests with Docker",c={unversionedId:"tooling/integration_tests",id:"tooling/integration_tests",isDocsHomePage:!1,title:"Integration tests with Docker",description:"Integration testing",source:"@site/docs/tooling/integration_tests.md",sourceDirName:"tooling",slug:"/tooling/integration_tests",permalink:"/docs/tooling/integration_tests",editUrl:"https://github.com/iotaledger/goshimmer/tree/develop/documentation/docs/tooling/integration_tests.md",version:"current",frontMatter:{},sidebar:"docs",previous:{title:"Docker private network",permalink:"/docs/tooling/docker_private_network"},next:{title:"How to do a release",permalink:"/docs/teamresources/release"}},p=[{value:"How to run",id:"how-to-run",children:[]},{value:"Creating tests",id:"creating-tests",children:[]},{value:"Other tips",id:"other-tips",children:[]}],u={toc:p};function d(e){var t=e.components,s=(0,o.Z)(e,a);return(0,i.kt)("wrapper",(0,r.Z)({},u,s,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"integration-tests-with-docker"},"Integration tests with Docker"),(0,i.kt)("p",null,(0,i.kt)("img",{alt:"Integration testing",src:n(2466).Z,title:"Integration testing"})),(0,i.kt)("p",null,"Running the integration tests spins up a ",(0,i.kt)("inlineCode",{parentName:"p"},"tester")," container within which every test can specify its own GoShimmer network with Docker as schematically shown in the figure above."),(0,i.kt)("p",null,"Peers can communicate freely within their Docker network and this is exactly how the tests are run using the ",(0,i.kt)("inlineCode",{parentName:"p"},"tester")," container.\nTest can be written in regular Go style while the framework provides convenience functions to create a new network, access a specific peer's web API or logs."),(0,i.kt)("h2",{id:"how-to-run"},"How to run"),(0,i.kt)("p",null,"Prerequisites: "),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"Docker 17.12.0+"),(0,i.kt)("li",{parentName:"ul"},"Docker compose: file format 3.5")),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre"},"# Mac & Linux\ncd tools/integration-tests\n./runTests.sh\n")),(0,i.kt)("p",null,"The tests produce ",(0,i.kt)("inlineCode",{parentName:"p"},"*.log")," files for every networks' peer in the ",(0,i.kt)("inlineCode",{parentName:"p"},"logs")," folder after every run."),(0,i.kt)("p",null,"On GitHub logs of every peer are stored as artifacts and can be downloaded for closer inspection once the job finishes."),(0,i.kt)("h2",{id:"creating-tests"},"Creating tests"),(0,i.kt)("p",null,"Tests can be written in regular Go style. Each tested component should reside in its own test file in ",(0,i.kt)("inlineCode",{parentName:"p"},"tools/integration-tests/tester/tests"),".\n",(0,i.kt)("inlineCode",{parentName:"p"},"main_test")," with its ",(0,i.kt)("inlineCode",{parentName:"p"},"TestMain")," function is executed before any test in the package and initializes the integration test framework."),(0,i.kt)("p",null,"Each test has to specify its network where the tests are run. This can be done via the framework at the beginning of a test."),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"// create a network with name 'testnetwork' with 6 peers and wait until every peer has at least 3 neighbors\nn := f.CreateNetwork(\"testnetwork\", 6, 3)\n// must be called to create log files and properly clean up\ndefer n.Shutdown() \n")),(0,i.kt)("h2",{id:"other-tips"},"Other tips"),(0,i.kt)("p",null,"Useful for development is to only execute the test you're currently building. For that matter, simply modify the ",(0,i.kt)("inlineCode",{parentName:"p"},"docker-compose.yml")," file as follows:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-yaml"},"entrypoint: go test ./tests -run <YOUR_TEST_NAME> -v -mod=readonly\n")))}d.isMDXComponent=!0},3905:function(e,t,n){"use strict";n.d(t,{Zo:function(){return p},kt:function(){return f}});var r=n(7294);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function a(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function s(e,t){if(null==e)return{};var n,r,o=function(e,t){if(null==e)return{};var n,r,o={},i=Object.keys(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var l=r.createContext({}),c=function(e){var t=r.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):a(a({},t),e)),n},p=function(e){var t=c(e.components);return r.createElement(l.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},d=r.forwardRef((function(e,t){var n=e.components,o=e.mdxType,i=e.originalType,l=e.parentName,p=s(e,["components","mdxType","originalType","parentName"]),d=c(n),f=o,m=d["".concat(l,".").concat(f)]||d[f]||u[f]||i;return n?r.createElement(m,a(a({ref:t},p),{},{components:n})):r.createElement(m,a({ref:t},p))}));function f(e,t){var n=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var i=n.length,a=new Array(i);a[0]=d;var s={};for(var l in t)hasOwnProperty.call(t,l)&&(s[l]=t[l]);s.originalType=e,s.mdxType="string"==typeof e?e:o,a[1]=s;for(var c=2;c<i;c++)a[c]=n[c];return r.createElement.apply(null,a)}return r.createElement.apply(null,n)}d.displayName="MDXCreateElement"},2466:function(e,t,n){"use strict";t.Z=n.p+"assets/images/integration-testing-a5a2fd4ebdfb3fb42cd75d867de81efd.png"}}]);