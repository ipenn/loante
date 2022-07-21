import gzip

def gzip_str(string_: str) -> bytes:
    return gzip.compress(string_.encode())


def gunzip_bytes_obj(bytes_obj: bytes) -> str:
    return gzip.decompress(bytes_obj).decode()


string_ = 'hello there!'

gzipped_bytes = gzip_str(string_)
gzipped_bytes = "\u001f\b\u0000\u0000\u0000\u0000\u0000\u0000\u0000½ÛnÛ¸\u0016_¥ðu:ð¢H\u001drÆm\u0011tO3Ó´\u001d`o\u000f\u0006´ÅØÚ%\r%%ñ\u0014y÷¡HZVJvA.¤Ejé#ù[?uàÉñå\u0016ÍU:9Àt\n³IÍø}¶d_·\u0015\u0003LÏ\u0026M¶auC7Ø\u000f\tD0$aÎ\u0026EY,E­\tÿÜn\u003eBÔ.Ò«w×á§jvÕþñáÛÕ÷Oü~ñ©¥]ÚlU\u0004â0ø\u0010ÏÞ¿8d\u0016]¾ÂÀ»Y\u0010Æ¢fJ\u001b*jþORÖ\\\u0015·å|r.ö\u0017yË²lÖjwC\u0017iÊY]ýù\u0004Ä\u001fOÎæn\få´ØÒùä©±æ¡äwêÈ¬ÅQü\u000bîÈÃ\u001e²Ûì²lF\u0014î¼\\Ò\u0026+\u000buÐªªÕF.M2#\u0007çe±ÚG$Q§\u0015\u0007iÆË¬Ùª#Ä\u003cd_ö\u001fÊWlÕ\rÏ¬\u0014Ç±J°¦\u003c} é\u0026)ËUk§Ó@fYpZ¤2DS\u0019P½öy×\u0007\u0019-oi±e\u0015/ÓVBõ\bÖÛºaï×ª­\u0002\u0011,á,g´fÏbuz7¬*vUæõ¶Î4¿ÉþQõ`ªê3Ñ\\\bdÁ¸j\u001aQ£´¬ÚNd2\u0026¶eL¶¯G\u001fìÉ(Ò5¿Ò;qúa7´hoé²i9ãýÐÓ´oß¢\u0014}¨\u0007k´\u001b\fÑÕr¸º®Ëê*§júm\u0011ÏÔá\u0018c5(uSrºÒcÂéækÙô-ï3É¢o5]äl¤lÃ6%ß^\n4]85ßd×OúAýåìRêS¶iV¾lºFw.\f^\u0015ÏÙFÀÖ(ÚG\u0007\tdô\u003eKÙ0AÐ\u0007\u00075e0-\u001f¼¤i\u0007Õ\u0007 ´ªò¬ÿÙüïô#ª·2*º¼»X©\u0002½­\n²âk\u0026\u000f\u0018\\nâH\u0014´¥à^)²?\u000füT*]pYÊ\u001f¦lámNW\u001d/(Z­Ax:\u001bg\u0005\u001b+¼\u001e+±ÂËY\u0015½\u001e+\u001acE/g\rl¬Áë±\u0006c¬ÁËY±\u0015¿\u001e+\u001ecÅ/g%6Vòz¬d¼5´±¯Ç\u001a±/gl¬Ñë±Fc¬ÑËYc\u001bküz¬ñ\u0018kürÖÄÆ¼\u001ek2Æ8Xÿìæ\u0027eÑP9Wê\\¬ZÎ\u0017\u0013qÆ\u0010#\u003dÁi+1ce\u001af\u003eASHÞÂôí\u0014Þ\u0000+ã¹®¶äÌQmº«\u0026\u0027ÛêÔ,ÕS\u003d@\u0005ûv\u0015ÊHe¿Ñ\u0027:°Ò\u000ey¢CV:ä \u000b\u003cÑ\u0005VºÀA\u003dÑa+\u001dvÐ\u0011OtÄJG\u001ct¡\u0027ºÐJ\u0017:è\"Ot.rÐÅèb+]ì K\u003cÑ%VºÄu5öe\u0016\u000e·\u0000§]xó\u000ba¸\u001c\u0003|YÃ3Àe\u001aàË5\u001c¶\u0001.ß\u0000_Æáp\u000epY\u0007øò\u000eyË\u003dÀ}8ü\u0003\\\u0006\u0002¾\u001cÄa!àò\u0010ðe\"\u000e\u0017\u0001/\u001fq\u0018\t¸\u0004ùr\u0012dw\u0012är\u0012äËIÝIóÞÃÛÍãîÃå$È »  _NìNÒ\u0015©\u001buçú.½Ø¿føù6½Ñ÷öúDúÎZ¥-ß\u003d±\u000e¦²Ñ£À\bD`O\u001c\bÙ\u0013\u0005D(°\u0027ÂDØHí#\u00111\u0012\u0011{¢Ð(4\u0012öD#Qd$ìbG¢ØH\u0014Û\u0013%D(q\bÒ%m0´]ÚvÛP7v¨\u001b\\ò\u0006CßØ¡op\t\u001c\fcÂÁ%q04\u001d\u001a\u0007ÈÁP9v¨\u001c\\2\u0007CçØ¡sp\t\u001d\f¥cÒÁ%u0´\u001dZ\u0007ØÁP;v¨\u001d¹Ô\fµ\u0013ÚKíÈP;q]Ë\u0017sCíÄ¡väR;2ÔN\u001cjG.µ#Cí\u0004K\u0027­7µ6ÒÁ5cÞ\u0014\u003eêÍEª÷äÝ~ÓßÐ¦í¦\u000f¾Ì(ö¯|û¹D\u0014\u0027b\u0002\u0012î\u003e\u0013(Òk2þ-Yv?- èf\u000bÞ\u0000\u0011zÞÏ\u003d_Æ\b\u0019ÁÊ\b~\u0018ÑaFdeD~\u0018Ã1ðÃ\u000f3b+#öÃH\u000e3\u0012+#ñÃ\u0018\u001ef\f­¡\u001fÆè0cdeü0Æ\u0019c+cì19ÌX\u0019\u0013O×ð#\u0006ìN\u0003¾¬æ\u0018¯q\u0027·#ì\u0006ì~\u0003\f\u0007p\u001c°[\u000exò\u001c8ÂtÀî:àÉvà\bß\u0001»ñ\u0027ç#¬\u0007ìÞ\u0003Ì\u0007p\u001f°Û\u000fxò\u001f8ÂÀî@àÉà\b\u000f\u0002»\t\u0027\u0017BG¸\u0010²»\u0010òäBè\b\u0017Bv\u0017B¾îy¹éqÜõxr!t\u000b!»\u000b!O.p!dw!t\u000bu\u000eÔ\u0007ê\u0017ù¢í!\u0014Ã§õJô\u001bz§BÝ\n­Y¶Z«\u0015\t~\u001cÿ¥rGw\u003e\u0015éWW|«G\u0018:¬Wb||V¹{\u0027q S¡\u0006\u0026\u001a\u0018h0\u0006£hp\u0012\u001a2Ñ\f44FÑÐIh\u0016h\u0016£\u0005£hÁIhØDÃ\u0026\u001a6Ðð8\u001a\u001eEÃ\u0027¡\u0011\u0013hÄ@#ãhd\u0014\u0016h¡\u0016\u001ahá8Z8\u0016\u0016h\u0016\u0019hÑ8Z4\u0016\u0016h±\u0016\u001bhñ8Z\u003c\u0016h\u0018hÉ8Z2\u001cÖ\u0019Ä\u0015Ó|F\u001b:ëôê\u003eõ¢w¿Ó\u0019Vê°ØÊTlEuhEuD¯è»®DÚ¦äûï½µ×þTÞw,ÌêoõbÆ\u0016íJÄoi^³Îö²Mëº»XW1+V¿ñòqû[É\u001b³è{U\u003c\u000b~)Ë¾\\ÁuËÅ¶ÒO»µ;{¤U¥×­¿?,JÎË\u0007ÝÒÝÎ_/Ü\u0013U¯ì}.\u000fUë\u0016cmæÞNCùÞ\u003dîÆF¿c\u001d,ÙºØèõó\tÑn:(ýÊøFM\u0018t\u001f~f\u000fr_/Ûª/Ú4SG÷ççå½^j÷,r\u0010*/iqÓ.ýêÎ£}\u0006U f,b,\u0006ü(*öXe|¤\u0013¢7(8\u0027øÕ8«èöyÎòñ´e3º­÷AUo·nQ·û¿kZü\rºâ\u001f\u0019¸)ëÁ@\rJ\u003ePÞ¬\u00197J-ç¬hV~Jûôg§¦Ý]Ue÷±ñ~\u0005Äî¹Zè¹« Ç²­÷ÇD»\u000f\u003eE¬ÿíÜ°ürúÿ7ïhNïxV¯\u000bZ¨F¦ÝªÄÁ2N\u003diÜ}\u0014,²üZ.²\\Ï(1Ä\u0004\u0010Æ}á\r{ì¦£Ý\u0002¼E\u0026\u001a?Û_1¦è-^÷U_ç\u001daÃ[ö4yú\u0017C÷²¤\u003c\u0000\u0000"
_bytes = bytes(gzipped_bytes, 'utf-8')
print(_bytes)
original_string = gunzip_bytes_obj(_bytes)
print(original_string)