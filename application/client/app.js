'use strict';

var app = angular.module('application', []);

app.controller('AppCtrl', function($scope, appFactory){
   $("#success_init").hide();
   $("#success_qurey").hide();
   $scope.initAB = function(){
       appFactory.initAB($scope.abstore, function(data){
           if(data == "")
           $scope.init_ab = "success";
           $("#success_init").show();
       });
   }
   $scope.queryAB = function(){
       appFactory.queryAB($scope.walletid, function(data){
           $scope.query_ab = data;
           $("#success_qurey").show();
       });
   }
});
app.factory('appFactory', function($http){
      
    var factory = {};
 
    factory.initAB = function(data, callback){
        $http.get('/init?a='+data.a+'&aval='+data.aval+'&b='+data.b+'&bval='+data.bval+'&c='+data.c+'&cval='+data.cval).success(function(output){
            callback(output)
        });
    }
    factory.queryAB = function(name, callback){
        $http.get('/query?name='+name).success(function(output){
            callback(output)
        });
    }
    return factory;
 });

 function isNumberKey(evt) {
    var charCode=(evt.which) ? evt.which : event.keyCode;
    if (charCode > 31 && (charCode < 48 || charCode > 57)) {
        return false;
    }
    return true;
 }

 function CheckPrice() {
    var inputElement = document.getElementById("X");
    var resultElement = document.getElementById("result");

    var inputValue = parseFloat(inputElement.value);
    if(!isNaN(inputValue)) {
        var increasedValue = inputValue * 1.03;
        resultElement.textContent = Math.floor(increasedValue);
    } else {
        alert("올바른 숫자를 입력하시오.")
    }
 }
 $(document).
        ready(function () {
            // Form submit handlers
            $('#initForm').submit(function (event) {
                event.preventDefault();
                var formData = $(this).serialize();
                $.get('/init?' + formData, function (data) {
                    $('#response').html(data);
                });
            });

            $('#invokeForm').submit(function (event) {
                event.preventDefault();
                var formData = $(this).serialize();
                $.get('/invoke?' + formData, function (data) {
                    
                    $('#response').html(data);
                });
            });

            $('#queryForm').submit(function (event) {
                event.preventDefault();
                var formData = $(this).serialize();
                $.get('/query?' + formData, function (data) {
                    $('#response').html(data);
                });
            });
        });