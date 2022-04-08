"use strict";
define([
  "text!temple/idig/index.html",
  "vue",
  "component/gl-btn/index",
  "component/gl-toggle-btn/index",
  "component/gl-input/index",
  "component/gl-tooltip/index",
], function (temp, Vue, gl_btn, gl_toggle, gl_input, gl_tooltip) {
  var vueComponent = Vue.extend({
    template: temp,
    data: function data() {
      return {
        enabled: false,
        name: "iDig Station",
        log: [],
      };
    },
    components: {
      "gl-btn": gl_btn,
      "gl-tg-btn": gl_toggle,
      "gl-input": gl_input,
      "gl-tooltip": gl_tooltip,
    },
    mounted: function mounted() {
      this.getStation();
    },
    methods: {
      toggleEnabled: function toggleEnabled() {
        this.enabled = !this.enabled;
      },
      getStation: function getStation() {
        var _this = this;
        _this.$store
          .dispatch("call", {api: "get_idig_station"})
          .then(function (result) {
            if (result.success) {
              _this.enabled = result.enabled;
              _this.name = result.name;
              _this.log = result.log;
            }
          });
      },
        setStation: function setStation() {
          var _this = this;
          _this.$store
            .dispatch("call", {
              api: "set_idig_station",
              data: { enabled: _this.enabled, name: _this.name },
            })
            .then(function (result) {
              if (result.success) {
                _this.$message({
                  type: "success",
                  api: "set_idig_station",
                  msg: result.code,
                });
                _this.enabled = result.enabled;
                _this.name = result.name;
                _this.log = result.log;
                } else {
                _this.$message({
                  type: "error",
                  api: "set_idig_station",
                  msg: result.code,
                });
              }
            });
        },
    },
    beforeRouteEnter: function beforeRouteEnter(to, from, next) {
      next(function (vm) {
        $("#router-visual").slideUp();
        if ($(".clsLink2" + vm.$route.path.split("/")[1]).hasClass("bar")) {
          $(".bar.active").removeClass("active");
          $(".clsLink2" + vm.$route.path.split("/")[1]).addClass("active");
          $("#applications").collapse("hide");
          $("#moresetting").collapse("hide");
          $("#system").collapse("hide");
          $("#vpn").collapse("hide");
        }
      });
    },
  });
  return vueComponent;
});
