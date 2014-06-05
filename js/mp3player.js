function Mp3Player($scope, $http) {
  $scope.files = [];
  $scope.loading = true;
  $scope.startingPath = "/mp3/";
  $scope.currentPath = [];
  $scope.error = false;
  $scope.playing = "";
  $scope.playingDirectory = [];

  $http({method: 'GET', url: $scope.startingPath}).
    success(function(data, status, headers, config) {
      console.log(data);
      $scope.files = data;
      $scope.loading = false;
    }).
    error(function(data, status, headers, config) {
      $scope.loading = false;
      $scope.error = true;
    });

  $scope.manageFile = function(file, index) {
    if(file.Directory) {
        if(file.Name == "..")
            $scope.currentPath.pop();
        else
            $scope.currentPath.push(file.Name);
        $http({method: 'GET', url: $scope.startingPath + $scope.currentPath.join('/')}).
            success(function(data, status, headers, config) {
              console.log(data);
              $scope.files = data;
              if($scope.currentPath.length > 0) $scope.files.unshift({Name: "..", Directory: true });
              $scope.loading = false;
            }).
            error(function(data, status, headers, config) {
              $scope.loading = false;
              $scope.error = true;
            });
    } else {
        $scope.play(file.Name, index);
    }
}

  $scope.play = function(name, index, fullPath) {
    $scope.playingDirectory = [];
    var currentPath = $scope.startingPath + $scope.currentPath.join('/') + "/";
    for(var i=index+1; i < $scope.files.length; i++) {
        if(!$scope.files[i].Directory)
            $scope.playingDirectory.push({full: currentPath + $scope.files[i].Name, name: $scope.files[i].Name});
    }
    var fullname = "";
    if(fullPath)
        fullname = fullPath;
    else
        fullname = currentPath + name;
    $('#audioplayer').attr('src', fullname);
    $scope.playing = name;
    if(fullPath) $scope.$digest();
  }

  $scope.stop = function(name) {
    if($scope.playing == name) {
      $('#audioplayer').attr('src', '');
      $scope.playing = '';
    }
  }

  mp3ended = function() {
    console.log('Ended!');
    if($scope.playingDirectory.length > 0) {
        var nextMp3 = $scope.playingDirectory.shift();
        $scope.play(nextMp3.name, undefined, nextMp3.full);
    } else {
        $scope.playing = '';
        $scope.$digest();
    }
  }

  getScope = function() {
    return $scope;
  }
}