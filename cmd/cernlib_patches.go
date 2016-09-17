package cmd

func prep_cernlib_patches() {
	text :=
		`--- Install_cernlib	2010-07-31 17:16:59.000000000 -0400
+++ Install_cernlib.patched	2016-09-10 10:09:09.414352008 -0400
@@ -1,4 +1,4 @@
-#!/bin/sh
+#!/bin/bash
 
 echo "===================="
 echo "CERNLIB installation"
@@ -14,6 +14,15 @@
 
 ./Install_cernlib_src
 
+# patch kstring.h
+pushd 2005/src/packlib/kuip/kuip
+patch < $CERN/kstring.h.patch
+popd
+
+# patch imake files
+pushd 2005/src/config
+patch < $CERN/Imake.cf.patch
+popd
 
 # Define the cernlib version
 
@@ -28,9 +37,8 @@
       end
 EOF
 GCCVSN=` + "`" + `cpp -dM comptest.F | grep __VERSION__ | cut -d" " -f3 | cut -c 2` + "`" +
			`
-FC=" "
+FC=gfortran
 [ "$GCCVSN" = "3" ]&&FC=g77
-[ "$GCCVSN" = "4" ]&&FC=gfortran
 if [ "$GCCVSN" = " " ]; then
   echo " "
   echo "====================================="`
	write_text("Install_cernlib.patch", text)

	text =
		`--- kstring.h	2015-03-25 20:15:56.852847033 -0400
+++ kstring.h.patched	2015-03-25 20:08:53.879725909 -0400
@@ -48,6 +48,10 @@
                      const char* str4 );
 extern char* str5dup( const char* str1, const char* str2, const char* str3,
                      const char* str4, const char* str5 );
+
+#ifdef strndup
+# undef strndup /* otherwise the next function declaration may bomb */
+#endif
 extern char* strndup( const char* buf, size_t n );
 extern char* stridup( int i );
 `
	write_text("kstring.h.patch", text)

	text =
		`--- Imake.cf	2006-06-15 04:59:55.000000000 -0400
+++ Imake.cf.patched	2016-09-10 09:31:25.836169871 -0400
@@ -393,6 +393,12 @@
 #elif __GNUC__ == 3
 #undef __GNUC__
 #define GCC3
+#elif __GNUC__ == 5
+#undef __GNUC__
+#define GCC4
+#elif __GNUC__ == 6
+#undef __GNUC__
+#define GCC4
 #else
 /*  old linux compiler suite (default)  */
 #undef __GNUC__`
	write_text("Imake.cf.patch", text)
}
