file {'/tmp/build-dir':
  path    => "/tmp/build-dir",
  ensure  => directory,
  mode    => 0640
}

define util::exec($command, $cwd = "/tmp", $unless = "false", 
  $onlyif = "true", $refreshonly = false) {
  exec { "$title":
    path        => '/bin:/usr/bin:/usr/local/sbin:usr/sbin:/sbin',
    environment => 'HOME=/root',
    command     => "$command",
    user        => 'root',
    group       => 'root',
    logoutput   => on_failure,
    unless      => $unless,
    onlyif      => $onlyif,
    cwd         => "$cwd",
    timeout     => 0,
    refreshonly => $refreshonly,
  }
}

define util::append($file = $title, $content) {
  util::exec{ "append::$content::$title":
    command => "echo $content >> $file",
    unless => "grep \"$content\" $file"
  }
}

define util::src-install($path, $tar_flags, $name, $ext, $config_flags = "", 
  $make_flags = "-j3") {
  notify {
    "Downloading and untarring $name from $path/$name.$ext...": 
    require     => File['/tmp/build-dir'],
  }
  ->
  util::exec {"get::$path/$name.$ext":
    cwd         => "/tmp/build-dir",
    command     => "wget $path/$name.$ext && tar $tar_flags $name.$ext",
    unless      => "which $title",
  }
  ->
  notify {
    "Building $name...": 
  }
  ->
  util::exec {"make::$path/$name.$ext":
    cwd         => "/tmp/build-dir/$name",
    command     => "sh configure $config_flags && make $make_flags \
      && make install",
    unless      => "which $title",
  }
}

define apt::ppa($owner = $title, $ppa) {
  util::exec {"apt::ppa $owner/$ppa":
    command     => "add-apt-repository ppa:$owner/$ppa",
    unless      => "ls /etc/apt/sources.list.d/ | grep $owner-$ppa",
    notify      => Util::Exec["touch::/tmp/apt-update-$owner-$ppa"],
  }
  ->
  util::exec {"touch::/tmp/apt-update-$owner-$ppa":
    command => "touch /tmp/apt-update-$owner-$ppa",
    refreshonly => true,
  }
}

define apt::update($force = 'false') {
  if $force == 'true' {
    util::exec {"apt::update::$title":
      command     => "apt-get update --fix-missing",
    }
  } else {
    util::exec {"apt::update::$title":
      command     => "apt-get update --fix-missing",
      onlyif      => "ls -l /tmp/ | grep apt-update",
    }
    util::exec {"rm::/tmp/apt-update-$title":
      command => "rm -f /tmp/apt-update-*",
      onlyif      => "ls -l /tmp/ | grep apt-update"
    }
  }
}

define apt::update-alternatives($command = $title, $oldpath, $newpath) {
  util::exec {"apt::update-alternatives $command :: $oldpath -> $newpath":
    command     => "update-alternatives --install $oldpath $command $newpath 20",
    unless      => "update-alternatives --get-selections | grep $command | grep $newpath"
  }
}

package { 'python-software-properties':
  ensure => installed
}

apt::ppa {'ubuntu-toolchain-r':
  ppa => 'test',
  require => Package['python-software-properties'],
}

apt::ppa {'tortoisehg-ppa':
  ppa => 'releases',
  require => Package['python-software-properties'],
}

apt::update {'init':
  require => [Apt::Ppa['ubuntu-toolchain-r'], Apt::Ppa['tortoisehg-ppa']]
}

package { 'gcc-4.8': ensure => installed, require => Apt::Update['init'] }
->
apt::update-alternatives { "gcc" :
  oldpath => "/usr/bin/gcc",
  newpath => "/usr/bin/gcc-4.8"
}

package { 'gdb': ensure => installed, require => Apt::Update-alternatives['gcc'] }
package { 'vim': ensure => installed, require => Apt::Update['init'] }
package { 'build-essential': ensure => installed, require => Apt::Update['init']}
package { 'mercurial': ensure => installed, require => Apt::Update['init']}
package { 'git': ensure => installed, require => Apt::Update['init']}

/* package { 'gccgo-4.8': ensure => installed, require => Package['gcc-4.8'] } */
/* -> */
/* apt::update-alternatives { "go" : */
/*   oldpath => "/usr/bin/gcc", */
/*   newpath => "/usr/bin/gccgo-4.8" */
/* } */

file {'/home/godir':
  path    => "/home/godir",
  ensure  => directory,
  mode    => 0640
}
->
notify {
  "Downloading Go...": 
  require     => [Apt::Update-alternatives['gcc'], 
    Package['mercurial'],
    Package['build-essential']]
}
->
util::exec {"get::https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz":
  cwd         => "/home/godir",
  command     => "wget https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz \
    && tar -zxvf go1.2.linux-amd64.tar.gz",
  unless      => "which go",
}
->
notify {
  "Building Go...": 
}
->
util::exec {"build::go":
  cwd         => "/home/godir/go/src",
  command     => "bash all.bash",
  unless      => "which go",
}
->
util::exec {"install::go":
  cwd         => "/home/godir/go/bin",
  command     => "ln -s /home/godir/go/bin/go /usr/local/bin/go \
    && ln -s /home/godir/go/bin/godoc /usr/local/bin/godoc \
    && ln -s /home/godir/go/bin/gofmt /usr/local/bin/gofmt",
  unless      => "which go",
}
