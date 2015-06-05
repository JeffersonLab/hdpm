#!/usr/bin/perl
open (VERSION, "gfortran --version |");
$line = <VERSION>;
split (/\s+/, $line);
$version = $_[3];
if ($version lt '4.2') {
    $patch = 'true';
} else {
    $patch = 'false';
}
close (VERSION);
print "$patch\n";
exit;
